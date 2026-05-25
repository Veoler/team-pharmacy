package repository

import (
	"github.com/Veoler/team-pharmacy/internal/models"
	"gorm.io/gorm"
)

type CategoryesRepository interface {
	GetAll() ([]models.Category, error)
	Create(category *models.Category) error
	GetAllSubcategoryes(categoryID uint) ([]models.Subcategory, error)
	CreateSubcategory(subcategory *models.Subcategory) error
	GetByID(id uint) (*models.Category, error)
	GetSubcategoryByID(id uint) (*models.Subcategory, error)
}

type gormCategoryesRepository struct {
	db *gorm.DB
}

func NewCategoryesRepository(db *gorm.DB) CategoryesRepository {
	return &gormCategoryesRepository{db: db}
}

func (r *gormCategoryesRepository) GetAll() ([]models.Category, error) {
	var categoryes []models.Category

	if err := r.db.Find(&categoryes).Error; err != nil {
		return nil, err
	}

	return categoryes, nil
}

func (r *gormCategoryesRepository) Create(category *models.Category) error {
	if category == nil {
		return nil
	}

	return r.db.Create(category).Error
}

func (r *gormCategoryesRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category

	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *gormCategoryesRepository) GetAllSubcategoryes(categoryID uint) ([]models.Subcategory, error) {
	var subcategoryes []models.Subcategory

	if err := r.db.Where("category_id = ?", categoryID).Find(&subcategoryes).Error; err != nil {
		return nil, err
	}

	return subcategoryes, nil
}

func (r *gormCategoryesRepository) CreateSubcategory(subcategory *models.Subcategory) error {
	if subcategory == nil {
		return nil
	}

	return r.db.Create(subcategory).Error
}

func (r *gormCategoryesRepository) GetSubcategoryByID(id uint) (*models.Subcategory, error) {
	var subcategory models.Subcategory

	if err := r.db.First(&subcategory, id).Error; err != nil {
		return nil, err
	}

	return &subcategory, nil
}
