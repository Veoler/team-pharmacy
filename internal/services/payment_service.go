package services

import (
	"errors"
	"time" 

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"gorm.io/gorm"
)

var ErrPaymentNotFound = errors.New("платеж не найден")

type PaymentService interface {
	CreatePayment(req models.PaymentCreateRequest) (*models.Payment, error)
	GetPaymentByID(id uint) (*models.Payment, error)
	GetPaymentFromOrder(orderID uint) ([]models.Payment, error)
	DeletePayment(id uint) error
}

type paymentService struct {
	payment repository.PaymentRepository
	order   repository.OrderRepository
}

func NewPaymentService(payment repository.PaymentRepository, order repository.OrderRepository) PaymentService {
	return &paymentService{payment: payment, order: order}
}

func (s *paymentService) CreatePayment(req models.PaymentCreateRequest) (*models.Payment, error) {
	if err := s.validatePaymentCreate(req); err != nil {
		return nil, err
	}

	order, err := s.order.GetByID(req.OrderID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrOrderNotFound
        }
        return nil, err
    }

	existingPayments, err := s.payment.GetFromOrder(req.OrderID)
	if err != nil {
		return nil, err
	}

	totalPaid := 0
	for _, p := range existingPayments {
		if p.Status == models.PayStatusSuccess {
			totalPaid += p.Amount
		}
	}

	if totalPaid+req.Amount > order.FinalPrice {
		return nil, errors.New("")
	}

	newPayment := &models.Payment{
		OrderID: req.OrderID,
		Amount: req.Amount,
		Method: req.Method,
		Status: req.Status,
		PaidAt: req.PaidAt,
	}

	if newPayment.Status == models.PayStatusSuccess {
		now := time.Now()
		newPayment.PaidAt = &now
	}

	if err := s.payment.Create(newPayment); err != nil {
		return nil, err
	}

	if newPayment.Status == models.PayStatusSuccess && (totalPaid+req.Amount) >= order.FinalPrice {
		order.Status = models.StatusPaid
		_ = s.order.UpdateStatusByID(order)
	}

	return newPayment, nil
}

func  (s *paymentService) GetPaymentByID(id uint) (*models.Payment, error) {
	payment, err := s.payment.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}

		return nil, err
	}

	return payment, nil
}

func (s *paymentService) GetPaymentFromOrder(orderID uint) ([]models.Payment, error) {
	_, err := s.payment.GetFromOrder(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
	}	

	return nil, err
	
}

func (s *paymentService) DeletePayment(id uint) error {
	if _, err := s.payment.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPaymentNotFound
		}

		return err
	}

	return s.payment.Delete(id)
}

func (s *paymentService) validatePaymentCreate(req models.PaymentCreateRequest) error {
	if req.Amount <= 0 {
		return errors.New("поле amount должно быть больше 0")
	}

	if req.OrderID <= 0 {
		return errors.New("поле order_id должно быть больше 0")
	}	

	if !isValidPayMethod(req.Method) {
		return errors.New("поле method должно быть одним из значений: card, cash, online_wallet")
	}

	if !isValidPayStatus(req.Status) {
		return errors.New("поле status должно быть одним из значений: pending, success, failed")
	}

	return nil
}

func isValidPayMethod(method models.PayMethod) bool {
	switch method {
	case models.PayMethodCard, 
	models.PayMethodCash, 
	models.PayMethodOnlineWallet:
		return true
	default: 
		return false
	}
}

func isValidPayStatus(status models.PayStatus) bool {
	switch status {
	case models.PayStatusPending,
	models.PayStatusSuccess,
	models.PayStatusFailed:
		return true
	default: 
		return false
	}
}
