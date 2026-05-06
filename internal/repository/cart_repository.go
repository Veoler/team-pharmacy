package repository

import (
	"github.com/Veoler/team-pharmacy/internal/models"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetCartByUserID(id uint) (*models.Cart, error)
	AddItem(cart *models.Cart) error
	GetItemByID(id uint) (*models.CartItem, error)
	AddQuantity(cart *models.Cart) error
	DeleteItem(id uint) error
	DeleteCart(id uint) error
}

type gormCartRepository struct {
	db *gorm.DB
}

func NewCartRepository(
	db *gorm.DB,
) CartRepository {
	return &gormCartRepository{db: db}
}

func (r *gormCartRepository) GetCartByUserID(id uint) (*models.Cart, error) {
	var cart models.Cart

	if err := r.db.Preload("Items").Where("user_id = ?", id).First(&cart).Error; err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r gormCartRepository) GetItemByID(id uint) (*models.CartItem, error) {
	var item models.CartItem

	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *gormCartRepository) AddItem(cart *models.Cart) error {
	if cart == nil {
		return nil
	}

	return r.db.Save(cart).Error
}

func (r *gormCartRepository) AddQuantity(cart *models.Cart) error {
	if cart == nil {
		return nil
	}

	return r.db.Save(cart).Error
}

func (r *gormCartRepository) DeleteItem(id uint) error {
	return r.db.Delete(&models.CartItem{}, id).Error
}

func (r *gormCartRepository) DeleteCart(id uint) error {
	return r.db.Where("user_id = ?", id).Delete(&models.Cart{}).Error
}
