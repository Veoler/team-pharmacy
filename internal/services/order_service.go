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
	GetAllByUserID(user_id uint) ([]models.OrdersInfo, error)
	UpdateOrderStatus(order_id uint, status models.OrderUpdateStatusRequest) (*models.Order, error)
}

type orderService struct {
	order     repository.OrderRepository
	cart      repository.CartRepository
	user      repository.UserRepository
	promocode repository.PromocodeRepository
}

func NewOrderService(
	order repository.OrderRepository,
	cart repository.CartRepository,
	user repository.UserRepository,
	promocode repository.PromocodeRepository,
) OrderService {
	return &orderService{order: order, cart: cart, user: user, promocode: promocode}
}

func (s *orderService) CreateOrder(req models.OrderCreateRequest) (*models.Order, error) {
	if req.UserID == nil {
		return nil, errors.New("user id не указан")
	}

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
			MedicineID:   *v.MedicineID,
			MedicineName: v.MedicineName,
			Quantity:     *v.Quantity,
			LineTotal:    v.LineTotal,
			PricePerUnit: *v.PricePerUnit,
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

	discauntTotal := 0
	if req.PromocodeID != nil {
		promocode, err := s.promocode.GetByID(*req.PromocodeID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrPromocodeNotFound
			}
			return nil, err
		}
		if promocode.IsActive {
			if promocode.DiscountType == models.DisTypePercent {
				discauntTotal = total * (promocode.DiscountValue / 100)
			} else if promocode.DiscountType == models.DisTypeFixed {
				discauntTotal = promocode.DiscountValue
			}
		}
		return nil, errors.New("этот промокод уже закончился")
	}

	final := total - discauntTotal
	// добавлено
	status := models.StatusDraft

	order := models.Order{
		UserID:          *req.UserID,
		Status:          status, // добавлено
		Items:           orderItem,
		TotalPrice:      total,
		FinalPrice:      final,
		DeliveryAddress: address,
		Comment:         comment,
		DiscountTotal:   discauntTotal,
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

func (s *orderService) GetAllByUserID(user_id uint) ([]models.OrdersInfo, error) {
	_, err := s.user.GetByID(user_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	orders, err := s.order.GetAllByUserID(user_id)
	if err != nil {
		return nil, err
	}

	ordersInfo := make([]models.OrdersInfo, 0, len(orders))

	for _, order := range orders {
		ordersInfo = append(ordersInfo, models.OrdersInfo{
			OrderID:    order.ID,
			FinalPrice: order.FinalPrice,
			Status:     order.Status,
		})
	}

	return ordersInfo, nil
}

func (s *orderService) UpdateOrderStatus(order_id uint, req models.OrderUpdateStatusRequest) (*models.Order, error) {
	order, err := s.order.GetByID(order_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	if req.Status == nil {
		return nil, errors.New("поле status обязательно")
	}

	var status models.Status
	status = *req.Status
	if !isValidStatus(status) {
		return nil, errors.New("такого статуса не существует")
	}
	order.Status = status

	if err := s.order.UpdateStatusByID(order); err != nil {
		return nil, err
	}

	return order, nil
}

func isValidStatus(status models.Status) bool {
	switch status {
	case models.StatusCanceled,
		models.StatusCompleted,
		models.StatusDraft,
		models.StatusPaid,
		models.StatusPendingPayment,
		models.StatusShipped:
		return true
	default:
		return false
	}
}
