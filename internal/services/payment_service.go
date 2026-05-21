package services

import (
	"errors"
	"time" 

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"gorm.io/gorm"
)

var ( 
	ErrPaymentNotFound = errors.New("платеж не найден")
	ErrPaymentExceedsTotal  = errors.New("сумма платежей превысит итоговую стоимость заказа")
)

type PaymentService interface {
	CreatePayment(req models.PaymentCreateRequest) (*models.Payment, *models.OrderPaymentSummary, error)
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

func (s *paymentService) CreatePayment(req models.PaymentCreateRequest) (*models.Payment, *models.OrderPaymentSummary, error) {
	if err := s.validatePaymentCreate(req); err != nil {
		return nil, nil, err
	}

	order, err := s.order.GetByID(req.OrderID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil, ErrOrderNotFound
        }
        return nil, nil, err
    }

	existingPayments, err := s.payment.GetFromOrder(req.OrderID)
	if err != nil {
		return nil,nil,  err
	}

	totalPaid := calcTotalPaid(existingPayments)


	if totalPaid+req.Amount > order.FinalPrice {
		return nil, nil, ErrPaymentExceedsTotal
	}

	newPayment := &models.Payment{
		OrderID: req.OrderID,
		Amount: req.Amount,
		Method: req.Method,
		Status: models.PayStatusPending,
	}

	if err := s.payment.Create(newPayment); err != nil {
		return nil, nil, err
	}

	if newPayment.Status == models.PayStatusSuccess {
		now := time.Now()
		newPayment.PaidAt = &now
 
		if totalPaid+req.Amount >= order.FinalPrice {
			order.Status = models.StatusPaid
			if err := s.order.UpdateStatusByID(order); err != nil {
				return nil, nil, err
			}
		}
	}
 
	summary := buildSummary(order, totalPaid+req.Amount)
 
	return newPayment, summary, nil
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
	if _, err := s.order.GetByID(orderID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	
	payments, err := s.payment.GetFromOrder(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}	

	return payments, nil
}

func (s *paymentService) DeletePayment(id uint) error {
	payment, err := s.payment.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPaymentNotFound
		}
		return err
	}
 
	if err := s.payment.Delete(id); err != nil {
		return err
	}

	if payment.Status == models.PayStatusSuccess {
		order, err := s.order.GetByID(payment.OrderID)
		if err != nil {
			return err
		}
 
		remainingPayments, err := s.payment.GetFromOrder(payment.OrderID)
		if err != nil {
			return err
		}
 
		totalPaid := calcTotalPaid(remainingPayments)
 
		if totalPaid < order.FinalPrice {
			order.Status = models.StatusPendingPayment
			if err := s.order.UpdateStatusByID(order); err != nil {
				return err
			}
		}
	}
 
	return nil
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

func calcTotalPaid(payments []models.Payment) int {
	total := 0
	for _, p := range payments {
		if p.Status == models.PayStatusSuccess {
			total += p.Amount
		}
	}
	return total
}

func buildSummary(order *models.Order, totalPaid int) *models.OrderPaymentSummary {
	var status models.PaymentStatus
	switch {
	case totalPaid == 0:
		status = models.PaymentStatusUnpaid
	case totalPaid < order.FinalPrice:
		status = models.PaymentStatusPartial
	default:
		status = models.PaymentStatusPaid
	}
 
	return &models.OrderPaymentSummary{
		OrderID:       order.ID,
		FinalPrice:    order.FinalPrice,
		PaidAmount:    totalPaid,
		PaymentStatus: status,
	}
}
