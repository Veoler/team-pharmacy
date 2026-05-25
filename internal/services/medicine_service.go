package services

import (
	"errors"
	"fmt"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"gorm.io/gorm"
)

var ErrMedicineNotFound error = errors.New("лекарство по айди не найдено")

type MedicineService interface {
	GetAll() ([]models.MedicineInfo, error)
	GetByID(id uint) (*models.Medicine, error)
	Create(req models.MedicineCreateRequest) (*models.Medicine, error)
	Update(id uint, req models.MedicineUpdateRequest) (*models.Medicine, error)
	Delete(id uint) error
}

type medicineService struct {
	medicine repository.MedicineRepository
	category repository.CategoryesRepository
}

func NewMedicineService(
	medicine repository.MedicineRepository,
	category repository.CategoryesRepository,
) MedicineService {
	return &medicineService{
		medicine: medicine,
		category: category,
	}
}

func (s *medicineService) GetAll() ([]models.MedicineInfo, error) {
	medicines, err := s.medicine.GetAll()
	if err != nil {
		return nil, err
	}

	mediicinesInfo := make([]models.MedicineInfo, 0, len(medicines))

	for _, m := range medicines {
		mediicinesInfo = append(mediicinesInfo, models.MedicineInfo{
			ID:        m.ID,
			Name:      m.Name,
			Price:     m.Price,
			InStock:   m.InStock,
			AvgRating: m.AvgRating,
		})
	}

	return mediicinesInfo, nil
}

func (s *medicineService) GetByID(id uint) (*models.Medicine, error) {
	mediicne, err := s.medicine.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMedicineNotFound
		}
		return nil, err
	}

	return mediicne, nil
}

func (s *medicineService) Create(req models.MedicineCreateRequest) (*models.Medicine, error) {
	if err := s.validCreate(req); err != nil {
		return nil, fmt.Errorf("Ошибка валидации при создании лекарства: %w", err)
	}

	newMedicine := &models.Medicine{
		Name:                 *req.Name,
		Description:          *req.Description,
		Price:                *req.Price,
		InStock:              *req.InStock,
		StockQuantity:        *req.StockQuantity,
		CategoryID:           *req.CategoryID,
		SubCategoryID:        *req.SubCategoryID,
		Manufacturer:         *req.Manufacturer,
		PrescriptionRequired: req.PrescriptionRequired,
	}

	if err := s.medicine.Create(newMedicine); err != nil {
		return nil, err
	}

	return newMedicine, nil
}

func (s *medicineService) validCreate(req models.MedicineCreateRequest) error {
	if req.Name == nil {
		return errors.New("name для лекарства обязателен")
	}
	if *req.Name == "" {
		return errors.New("name не должен быть пустым")
	}

	if req.Description == nil {
		return errors.New("description для лекарства обязателен")
	}
	if *req.Description == "" {
		return errors.New("description не должен быть пустым")
	}

	if req.Price == nil {
		return errors.New("price для лекарства обязателен")
	}
	if *req.Price <= 0 {
		return errors.New("price должен быть больше 0")
	}

	if req.InStock == nil {
		return errors.New("in_stock для лекарства обязателен")
	}
	if *req.InStock && *req.StockQuantity == 0 {
		return errors.New("вы указали in_stock true но stock_quantity указали 0 это как")
	}

	if req.CategoryID == nil {
		return errors.New("category_id для лекарства обязателен")
	}
	if _, err := s.category.GetByID(*req.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
	}

	if req.SubCategoryID == nil {
		return errors.New("subcategory_id для лекарства обязателен")
	}
	if _, err := s.category.GetSubcategoryByID(*req.SubCategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSubCategoryNotFound
		}
	}

	if req.Manufacturer == nil {
		return errors.New("manufacturer для лекарства обязателен")
	}
	if *req.Manufacturer == "" {
		return errors.New("manufacturer для лекарства не может быть пустым")
	}

	return nil
}

func (s *medicineService) Update(id uint, req models.MedicineUpdateRequest) (*models.Medicine, error) {
	medicine, err := s.medicine.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMedicineNotFound
		}
		return nil, err
	}

	if req.Name != nil {
		medicine.Name = *req.Name
	}
	if req.Description != nil {
		medicine.Description = *req.Description
	}
	if req.Price != nil {
		medicine.Price = *req.Price
	}
	if req.InStock != nil {
		medicine.InStock = *req.InStock
	}
	if req.StockQuantity != nil {
		medicine.StockQuantity = *req.StockQuantity
	}
	if !*req.InStock {
		medicine.StockQuantity = 0
	}
	if req.CategoryID != nil {
		if _, err := s.category.GetByID(*req.CategoryID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrCategoryNotFound
			}
		}
		medicine.CategoryID = *req.CategoryID
	}
	if req.SubCategoryID != nil {
		if _, err := s.category.GetSubcategoryByID(*req.SubCategoryID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrSubCategoryNotFound
			}
		}
		medicine.SubCategoryID = *req.SubCategoryID
	}
	if req.Manufacturer != nil {
		medicine.Manufacturer = *req.Manufacturer
	}
	if req.PrescriptionRequired != nil {
		medicine.PrescriptionRequired = *req.PrescriptionRequired
	}

	if err := s.medicine.Update(medicine); err != nil {
		return nil, err
	}

	return medicine, nil
}

func (s *medicineService) Delete(id uint) error {
	if _, err := s.medicine.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMedicineNotFound
		}
		return err
	}

	return s.medicine.Delete(id)
}
