package services

import (
	"errors"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"gorm.io/gorm"
	// "gorm.io/gorm"
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
}

func NewPaymentRepository(payment repository.PaymentRepository) PaymentService {
	return &paymentService{payment: payment}
}

func (s *paymentService) CreatePayment(req models.PaymentCreateRequest) (*models.Payment, error) {
	if err := s.validatePaymentCreate(req); err != nil {
		return nil, err
	}

	// if err := s.ensureOrderExists(req.OrderID); err != nil {
		// return nil, err
	// }

	newPayment := &models.Payment{
		Amount: req.Amount,
		Method: req.Method,
		Status: req.Status,
		PaidAt: req.PaidAt,
	}

	if err := s.payment.Create(newPayment); err != nil {
		return nil, err
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
	if err := s.ensureOrderExists(orderID); err != nil {
		return nil, err
	}

	return s.payment.GetFromOrder(orderID)
}

func (s *paymentService) DeletePayment(id uint) error {
	// if _, err := s.payment.GetByID(id); err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
			// return ErrPaymentNotFound
		// }
// 
		// return err
	// }

	if _, err := s.GetPaymentByID(id); err != nil {
		return err
	}

	return s.payment.Delete(id)
}

func (s *paymentService) validatePaymentCreate(req models.PaymentCreateRequest) error {
	if req.Amount <= 0 {
		return errors.New("поле amount должно быть больше 0")
	}

	// if req.OrderID <= 0 {
		// return errors.New("поле order_id должно быть больше 0")
	// }	

	if !isValidPayMethod(req.Method) {
		return errors.New("поле method должно быть одним из значений: card, cash, online_wallet")
	}

	if !isValidPayStatus(req.Status) {
		return errors.New("поле status должно быть одним из значений: pending, success, failed")
	}

	return nil
}


func (s *paymentService) ensureOrderExists(orderID uint) error {
	if _, err := s.payment.GetByID(orderID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err // ErrOrderNotFound
		}

		return err
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
