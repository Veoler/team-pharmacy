package repository

import (
	"github.com/Veoler/team-pharmacy/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetOrdersByUserID(id uint) ([]models.Order, error)
	GetCartByUserID(id uint) (*models.Cart, error)
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(
	db *gorm.DB,
) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(user *models.User) error {
	if user == nil {
		return nil
	}

	return r.db.Create(user).Error
}

func (r *gormUserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *gormUserRepository) GetOrdersByUserID(id uint) ([]models.Order, error) {
	var orders []models.Order

	if err := r.db.Preload("Items").Where("user_id = ?", id).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *gormUserRepository) GetCartByUserID(id uint) (*models.Cart, error) {
	var cart models.Cart

	if err := r.db.Preload("Items").Where("user_id = ?", id).First(&cart).Error; err != nil {
		return nil, err
	}

	return &cart, nil
}
