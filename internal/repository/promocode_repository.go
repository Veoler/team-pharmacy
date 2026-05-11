package repository

import (
	"gorm.io/gorm"
	"github.com/Veoler/team-pharmacy/internal/models"
)

type PromocodeRepository interface {
	Create(promocode *models.Promocode) error
	GetAll() ([]models.Promocode, error)
	GetByID(id uint) (*models.Promocode, error)
	Update(promocode *models.Promocode) error
	Delete(id uint) error
	Validate(promocode *models.Promocode) (*models.Promocode, error)
}

type promocodeRepository struct {
	db *gorm.DB
}

func NewPromocodeRepository(db *gorm.DB) PromocodeRepository {
	return &promocodeRepository{db: db}
}

func(r *promocodeRepository) Create(promocode *models.Promocode) error {
	if promocode == nil {
		return nil
	}

	return r.db.Create(promocode).Error
}

func(r *promocodeRepository) GetAll() ([]models.Promocode, error) {
	var promocodes []models.Promocode

	if err := r.db.Find(&promocodes).Error; err != nil {
		return nil, err
	}

	return promocodes, nil
}

func(r *promocodeRepository) GetByID(id uint) (*models.Promocode, error) {
	var promocode models.Promocode

	if err := r.db.First(&promocode, id).Error; err != nil {
		return nil, err
	}

	return &promocode, nil
}

func(r *promocodeRepository) Update(promocode *models.Promocode) error {
	if promocode == nil {
		return nil
	}

	return r.db.Save(promocode).Error
}

func(r *promocodeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Promocode{}, id).Error
}

func(r *promocodeRepository) Validate(promocode *models.Promocode) (*models.Promocode, error) {
	return nil, nil
}


