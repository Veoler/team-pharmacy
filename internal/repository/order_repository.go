package repository

import (
	"github.com/Veoler/team-pharmacy/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetByID(id uint) (*models.Order, error)
	GetAllByUserID(id uint) ([]models.Order, error)
	UpdateStatusByID(order *models.Order) error
}

type gormOrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(
	db *gorm.DB,
) OrderRepository {
	return &gormOrderRepository{db: db}
}

func (r *gormOrderRepository) CreateOrder(order *models.Order) error {
	if order == nil {
		return nil
	}

	return r.db.Create(order).Error
}

func (r *gormOrderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order

	if err := r.db.Preload("Items").First(&order, id).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *gormOrderRepository) GetAllByUserID(id uint) ([]models.Order, error) {
	var orders []models.Order

	if err := r.db.Preload("Items").Where("user_id = ?", id).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *gormOrderRepository) UpdateStatusByID(order *models.Order) error {
	if order == nil {
		return nil
	}

	return r.db.Save(order).Error
}
