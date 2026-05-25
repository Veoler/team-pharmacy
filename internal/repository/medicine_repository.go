package repository

import (
	"github.com/Veoler/team-pharmacy/internal/models"
	"gorm.io/gorm"
)

type MedicineRepository interface {
	GetAll() ([]models.Medicine, error)
	GetByID(id uint) (*models.Medicine, error)
	Create(medicine *models.Medicine) error
	Update(medicine *models.Medicine) error
	Delete(id uint) error
}

type gormMedicineRepository struct {
	db *gorm.DB
}

func NewMedicineRepository(db *gorm.DB) MedicineRepository {
	return &gormMedicineRepository{
		db: db,
	}
}

func (r *gormMedicineRepository) GetAll() ([]models.Medicine, error) {
	var medicines []models.Medicine

	if err := r.db.Find(&medicines).Error; err != nil {
		return nil, err
	}

	return medicines, nil
}

func (r *gormMedicineRepository) GetByID(id uint) (*models.Medicine, error) {
	var medicine models.Medicine

	if err := r.db.First(&medicine, id).Error; err != nil {
		return nil, err
	}

	return &medicine, nil
}

func (r *gormMedicineRepository) Create(medicine *models.Medicine) error {
	if medicine == nil {
		return nil
	}

	return r.db.Create(medicine).Error
}

func (r *gormMedicineRepository) Update(medicine *models.Medicine) error {
	if medicine == nil {
		return nil
	}

	return r.db.Save(medicine).Error
}

func (r *gormMedicineRepository) Delete(id uint) error {
	return r.db.Delete(&models.Medicine{}, id).Error
}
