package services

import (
	"errors"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"gorm.io/gorm"
)

var ErrCartNotFound error = errors.New("cart по такому user_id не найден")
var ErrCartItemNotFound error = errors.New("cart_item по такому user_id не найден")

type CartService interface {
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

func (s *cartService) GetCartByUserID(id uint) (*models.Cart, error) {
	_, err := s.user.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	cart, err := s.cart.GetCartByUserID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartNotFound
		}
		return nil, err
	}

	return cart, nil
}

func (s *cartService) AddItem(req models.CartCreateUpdateRequest) (*models.Cart, error) {
	_, err := s.user.GetByID(req.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	cart, err := s.cart.GetCartByUserID(req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartNotFound
		}
		return nil, err
	}

	for _, newItem := range req.Items {
		found := false
		for i := range cart.Items {
			if cart.Items[i].MedicineID == newItem.MedicineID {
				cart.Items[i].Quantity += newItem.Quantity
				cart.Items[i].LineTotal = cart.Items[i].Quantity * cart.Items[i].PricePerUnit
				found = true
				break
			}
			if !found {
				cart.Items = append(cart.Items, newItem)
			}
		}
	}

	if err := s.cart.AddItem(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *cartService) AddQuantity(req models.CartCreateUpdateRequest) (*models.Cart, error) {
	_, err := s.user.GetByID(req.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	cart, err := s.cart.GetCartByUserID(req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartNotFound
		}
		return nil, err
	}

	for _, newItem := range req.Items {
		for i := range cart.Items {
			if cart.Items[i].MedicineID == newItem.MedicineID {
				cart.Items[i].Quantity += newItem.Quantity
				cart.Items[i].LineTotal = cart.Items[i].Quantity * cart.Items[i].PricePerUnit
				break
			}
		}
	}

	if err := s.cart.AddQuantity(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *cartService) DeleteItem(id uint) error {
	_, err := s.cart.GetItemByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrCartItemNotFound
	}

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

	if err := s.cart.DeleteCart(id); err != nil {
		return err
	}

	return nil
}
