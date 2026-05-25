package services

import (
	"errors"
	"fmt"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"gorm.io/gorm"
)

var ErrCartNotFound error = errors.New("cart по такому user_id не найден")
var ErrCartItemNotFound error = errors.New("cart_item по такому id не найден")

type CartService interface {
	AddItem(req models.CartCreateUpdateRequest) (*models.Cart, error)
	AddQuantity(userID, itemID *uint, newQuantity int) (*models.Cart, error)
	DeleteItem(id uint, userID uint) error
	DeleteCart(id uint) error
}

type cartService struct {
	cart repository.CartRepository
	user repository.UserRepository
}

func NewCartService(
	cart repository.CartRepository,
	user repository.UserRepository,
) CartService {
	return &cartService{cart: cart, user: user}
}

func (s *cartService) AddItem(req models.CartCreateUpdateRequest) (*models.Cart, error) {
	if err := s.validateAddItem(req); err != nil {
		return nil, fmt.Errorf("при валидации создания позиции для корзины возникла ошибка: %w", err)
	}

	_, err := s.user.GetByID(*req.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	cart, err := s.cart.GetCartByUserID(*req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = &models.Cart{
				UserID: *req.UserID,
				Items:  []models.CartItem{},
			}
		} else {
			return nil, err
		}
	}

	for _, newItem := range req.Items {
		found := false
		for i := range cart.Items {
			if *cart.Items[i].MedicineID == *newItem.MedicineID {
				if newItem.Quantity != nil {
					// *cart.Items[i].Quantity += *newItem.Quantity
					newQty := *cart.Items[i].Quantity + *newItem.Quantity
					cart.Items[i].Quantity = &newQty
					cart.Items[i].LineTotal = *cart.Items[i].Quantity * *cart.Items[i].PricePerUnit
				}
				found = true
				break
			}
		}
		if !found {
			newItem.ID = 0
			if newItem.Quantity != nil && newItem.PricePerUnit != nil {
				newItem.LineTotal = *newItem.Quantity * *newItem.PricePerUnit
			}
			cart.Items = append(cart.Items, newItem)
		}
	}

	if err := s.cart.AddItem(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *cartService) validateAddItem(req models.CartCreateUpdateRequest) error {
	if req.UserID == nil || *req.UserID == 0 {
		return errors.New("не указан ID пользователя")
	}

	if len(req.Items) == 0 {
		return errors.New("корзина не может быть пустой, добавьте хотя бы один товар")
	}

	for _, item := range req.Items {
		if item.MedicineID == nil {
			return errors.New("medicine_id для items обязателен")
		}
		if item.Quantity == nil {
			return errors.New("quantity для items обязателен")
		}
		if *item.Quantity <= 0 {
			return errors.New("количество для товара должно быть больше 0")
		}
		if item.PricePerUnit == nil {
			return errors.New("price_per_unit для items обязателен")
		}
		if *item.PricePerUnit <= 0 {
			return errors.New("цена для товара не может быть нулевой")
		}
	}

	return nil
}

func (s *cartService) AddQuantity(userID, itemID *uint, newQuantity int) (*models.Cart, error) {
	if userID == nil || *userID == 0 {
		return nil, errors.New("не указан ID пользователя")
	}

	_, err := s.user.GetByID(*userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	cart, err := s.cart.GetCartByUserID(*userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartNotFound
		}
		return nil, err
	}

	var targetItem *models.CartItem

	for i := range cart.Items {
		if cart.Items[i].ID == *itemID {
			targetItem = &cart.Items[i]
			break
		}
	}

	if targetItem == nil {
		return nil, ErrCartItemNotFound
	}

	targetItem.Quantity = &newQuantity
	targetItem.LineTotal = newQuantity * *targetItem.PricePerUnit

	if err := s.cart.AddQuantity(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *cartService) DeleteItem(id uint, userID uint) error {
	_, err := s.cart.GetItemByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrCartItemNotFound
	}

	if _, err := s.user.GetByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
	}
//добавлено
	if _, err := s.cart.GetCartByUserID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartNotFound
		}
	}
// до сюда
	if err := s.cart.DeleteItem(id); err != nil {
		return err
	}

	return nil
}

func (s *cartService) DeleteCart(id uint) error {
	_, err := s.user.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}
//добавлено
	if _, err := s.cart.GetCartByUserID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartNotFound
		}
	}
//до сюда
	if err := s.cart.DeleteCart(id); err != nil {
		return err
	}

	return nil
}
