package services

import (
	"errors"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"gorm.io/gorm"
)

var ErrOrdersNotFound error = errors.New("orders по такому user_id не найдены")
var ErrOrderNotFound error = errors.New("order по такому id не найден")

type OrderService interface {
	CreateOrder(req models.OrderCreateRequest) (*models.Order, error)
	GetByID(id uint) (*models.Order, error)
	GetAllByUserID(user_id uint) ([]models.Order, error)
	UpdateOrderStatus(order_id uint, status models.Status) (*models.Order, error)
}

type orderService struct {
	order repository.OrderRepository
	cart  repository.CartRepository
	user  repository.UserRepository
}

func NewOrderService(
	order repository.OrderRepository,
	cart repository.CartRepository,
	user repository.UserRepository,
) OrderService {
	return &orderService{order: order, cart: cart, user: user}
}

func (s *orderService) CreateOrder(req models.OrderCreateRequest) (*models.Order, error) {
	user, err := s.user.GetByID(*req.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	cart, err := s.cart.GetCartByUserID(*req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartNotFound
		}
		return nil, err
	}

	var total int
	orderItem := []models.OrderItem{}

	for _, v := range cart.Items {
		tempItem := models.OrderItem{
			MedicineID:   v.MedicineID,
			MedicineName: v.MedicineName,
			Quantity:     v.Quantity,
			LineTotal:    v.LineTotal,
			PricePerUnit: v.PricePerUnit,
		}

		orderItem = append(orderItem, tempItem)

		total += v.LineTotal
	}

	address := req.DeliveryAddress
	if req.DeliveryAddress == "" {
		address = user.DefaultAddress
	}
	comment := "нет коммента"
	if req.Comment != "" {
		comment = req.Comment
	}

	order := models.Order{
		UserID:          *req.UserID,
		Status:          models.StatusDraft,
		Items:           orderItem,
		TotalPrice:      total,
		FinalPrice:      total,
		DeliveryAddress: address,
		Comment:         comment,
	}

	if err := s.order.CreateOrder(&order); err != nil {
		return nil, err
	}

	if err := s.cart.DeleteCart(*req.UserID); err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *orderService) GetByID(id uint) (*models.Order, error) {
	order, err := s.order.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetAllByUserID(user_id uint) ([]models.Order, error) {
	_, err := s.user.GetByID(user_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	orders, err := s.order.GetAllByUserID(user_id)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *orderService) UpdateOrderStatus(order_id uint, status models.Status) (*models.Order, error) {
	order, err := s.order.GetByID(order_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	order.Status = status

	if err := s.order.UpdateStatusByID(order); err != nil {
		return nil, err
	}

	return order, nil
}
