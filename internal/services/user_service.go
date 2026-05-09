package services

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"gorm.io/gorm"
)

var ErrUserNotFound error = errors.New("user по такому id не найден")

type UserService interface {
	Create(req models.UserCreateRequest) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetOrdersByUserID(id uint) ([]models.Order, error)
	GetCartByUserID(id uint) (*models.Cart, error)
}

type userService struct {
	user repository.UserRepository
}

func NewUserService(
	user repository.UserRepository,
) UserService {
	return &userService{user: user}
}

func (s *userService) Create(req models.UserCreateRequest) (*models.User, error) {
	if err := s.isValidToCreate(req); err != nil {
		return nil, fmt.Errorf("При валидации для создания пользователя произошла ошибка: %w", err)
	}

	user := &models.User{
		FullName:       *req.FullName,
		Email:          *req.Email,
		Phone:          req.Phone,
		DefaultAddress: req.DefaultAddress,
	}

	if err := s.user.Create(user); err != nil {
		return nil, fmt.Errorf("Не удалось создать пользователя ошибка: %w", err)
	}

	return user, nil
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	user, err := s.user.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) GetOrdersByUserID(id uint) ([]models.Order, error) {
	_, err := s.user.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	orders, err := s.user.GetOrdersByUserID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrdersNotFound
		}
		return nil, err
	}

	return orders, nil
}

func (s *userService) GetCartByUserID(id uint) (*models.Cart, error) {
	_, err := s.user.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	cart, err := s.user.GetCartByUserID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartNotFound
		}
		return nil, err
	}

	return cart, nil
}

func (s *userService) isValidToCreate(req models.UserCreateRequest) error {
	if req.FullName == nil {
		return errors.New("full_name поле обязателен")
	}
	if *req.FullName == "" {
		return errors.New("full_name не должен быть пустым")
	}
	if req.Email == nil {
		return errors.New("e-mail поле обязателен")
	}
	if *req.Email == "" {
		return errors.New("e-mail не должен быть пустым")
	}
	addr, err := mail.ParseAddress(*req.Email)
	if err != nil || addr.Address != *req.Email {
		return errors.New(`e-mail не корректный должно быть примерно так "логин@example.com"`)
	}
	return nil
}
