package repository

import (
	"github.com/Veoler/team-pharmacy/internal/models"
	"gorm.io/gorm"
)


// После комментирования полей в models, здесь тоже
// закомментируйте поля которыe пока красным у вас горят

type PaymentRepository interface {
	Create(payment *models.Payment) error
	GetByID(id uint) (*models.Payment, error)
	GetFromOrder(orderID uint) ([]models.Payment, error)
	Delete(id uint) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func(r *paymentRepository) Create(payment *models.Payment) error {
	if payment == nil {
		return nil
	}

	return r.db.Create(&payment).Error
}

func(r *paymentRepository) GetByID(id uint) (*models.Payment, error) {
	var payment models.Payment

	if err := r.db.First(&payment, id).Error; err != nil {
		return nil, err
	}

	return &payment, nil
}

func(r *paymentRepository) GetFromOrder(orderID uint) ([]models.Payment, error) {
	var payments []models.Payment

	if err := r.db.
	Model(&models.Payment{}).
	Where("order_id = &", orderID).
	Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}

func(r *paymentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Payment{}, id).Error
}