package services

import (
	"errors"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"gorm.io/gorm"
)

var ErrCategoryNotFound error = errors.New("категория по айди не найдена")
var ErrSubCategoryNotFound error = errors.New("подкатегория по айди не найдена")

type CategoryesService interface {
	GetAll() ([]models.Category, error)
	Create(req models.CategoryCreateRequest) (*models.Category, error)
	GetAllSubcategoryes(categoryID uint) ([]models.Subcategory, error)
	CreateSubcategory(req models.SubcategoryCreateRequest) (*models.Subcategory, error)
}

type categoryesService struct {
	categoryes repository.CategoryesRepository
}

func NewCategoryesService(
	categoryes repository.CategoryesRepository,
) CategoryesService {
	return &categoryesService{
		categoryes: categoryes,
	}
}

func (s *categoryesService) GetAll() ([]models.Category, error) {
	categoryes, err := s.categoryes.GetAll()
	if err != nil {
		return nil, err
	}

	return categoryes, nil
}

func (s *categoryesService) Create(req models.CategoryCreateRequest) (*models.Category, error) {
	if req.Name == nil {
		return nil, errors.New("name для категории обязателен")
	}
	if *req.Name == "" {
		return nil, errors.New("name для категории не может быть пустым")
	}

	category := &models.Category{
		Name: *req.Name,
	}

	if err := s.categoryes.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryesService) GetAllSubcategoryes(categoryID uint) ([]models.Subcategory, error) {
	if _, err := s.categoryes.GetByID(categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
	}

	subcategoryes, err := s.categoryes.GetAllSubcategoryes(categoryID)
	if err != nil {
		return nil, err
	}

	return subcategoryes, nil
}

func (s *categoryesService) CreateSubcategory(req models.SubcategoryCreateRequest) (*models.Subcategory, error) {
	if req.CategoryID == nil {
		return nil, errors.New("category_id обязателен для подкатегории")
	}
	if _, err := s.categoryes.GetByID(*req.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
	}

	if req.Name == nil {
		return nil, errors.New("name обязателен для подкатегории")
	}
	if *req.Name == "" {
		return nil, errors.New("name не может быть пустой для подкатегории")
	}

	subcategory := &models.Subcategory{
		Name:       *req.Name,
		CategoryID: *req.CategoryID,
	}

	if err := s.categoryes.CreateSubcategory(subcategory); err != nil {
		return nil, err
	}

	return subcategory, nil
}
