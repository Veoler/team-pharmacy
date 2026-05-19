package services

import (
	"errors"
	"strings"

	"github.com/Veoler/team-pharmacy/internal/models"
	"github.com/Veoler/team-pharmacy/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrReviewNotFound = errors.New("отзыв не найден")
	ErrNotPurchased = errors.New("нельзя оставить отзыв: вы не заказывали это лекарство")
)

type MedicineGetter interface {
	GetByID(id uint) (*models.Medicine, error)
}

type ReviewService interface {	
	CreateReview(req *models.ReviewCreateRequest) (*models.Review, error)
	GetReviewsFromMedicine(medicineID uint) ([]models.Review, error)
	UpdateReview(id uint, req models.ReviewUpdateRequest) (*models.Review, error)
	DeleteReview(id uint) error
}

type reviewService struct {
	review		repository.ReviewRepository
	order		repository.OrderRepository
	medicine	MedicineGetter
}

func NewReviewService(
	review repository.ReviewRepository, 
	order	repository.OrderRepository, 
	medicine MedicineGetter,
) ReviewService {
	return &reviewService{review: review, order: order, medicine: medicine}
}

func (s *reviewService) CreateReview(req *models.ReviewCreateRequest) (*models.Review, error) {
	if err := s.validateReviewCreate(req); err != nil {
		return nil, err
	}

	if err := s.checkMedicineExists(req.MedicineID); err != nil {
		return nil, err
	}

	purchased, err := s.hasPurchasedMedicine(req.UserID, req.MedicineID)
	if err != nil {
		return nil, err
	}
	if !purchased {
		return nil, ErrNotPurchased
	}

	review := &models.Review{
		UserID:		req.UserID,
		MedicineID:	req.MedicineID,
		Rating:		req.Rating,
		Text:		req.Text,
	}

	if err := s.review.Create(review); err != nil {
		return nil, err
	}

	if err := s.updateAvgRating(req.MedicineID); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *reviewService) GetReviewsFromMedicine(medicineID uint) ([]models.Review, error) {
	if err := s.checkMedicineExists(medicineID); err != nil {
		return nil, err
	}
	
	return s.review.GetFromMedicine(medicineID)
}

func (s *reviewService) UpdateReview(id uint, req models.ReviewUpdateRequest) (*models.Review, error) {
	review, err := s.review.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}

	if err := s.applyReviewUpdate(review, req); err != nil {
		return nil, err
	}

	if err := s.review.Update(review); err != nil {
		return nil, err
	}

	if err := s.updateAvgRating(review.MedicineID); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *reviewService) DeleteReview(id uint) error {
	review, err := s.review.GetByID(id); 
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReviewNotFound
		}
		return err
	}

	if err := s.review.Delete(id); err != nil {
		return err
	}
 
	if err := s.updateAvgRating(review.MedicineID); err != nil {
		return err
	}
 
	return nil
}

func (s *reviewService) validateReviewCreate(req *models.ReviewCreateRequest) error {
	if req.UserID <= 0 {
		return errors.New("поле user_id должно быть больше 0")
	}

	if req.MedicineID <= 0 {
		return errors.New("поле medicine_id должно быть больше 0")
	}

	if req.Rating < 1 || req.Rating > 5 {
		return errors.New("поле rating должно быть в диапозоне от 1 до 5")
	}

	if strings.TrimSpace(req.Text) == "" {
		return errors.New("поле text не должно быть пустым")
	}

	return nil
}

func (s *reviewService) applyReviewUpdate(review *models.Review, req models.ReviewUpdateRequest) error {
	if req.Rating != nil {
		if *req.Rating < 1 || *req.Rating > 5{
			return errors.New("поле rating должно быть в диапозоне от 1 до 5")
		}
	review.Rating = *req.Rating
	}

	if req.Text != nil {
		trimmed := strings.TrimSpace(*req.Text)
		if trimmed == "" {
			return errors.New("поле text не должно быть пустым")
		}
		review.Text = trimmed
	}

	return nil
}

func (s *reviewService) hasPurchasedMedicine(userID uint, medicineID uint) (bool, error) {
	orders, err := s.order.GetAllByUserID(userID)
	if err != nil {
		return false, err
	}
 
	for _, order := range orders {
		if order.Status != models.StatusCompleted {
			continue
		}
		for _, item := range order.Items {
			if item.MedicineID == medicineID {
				return true, nil
			}
		}
	}
 
	return false, nil
}

func (s *reviewService) updateAvgRating(medicineID uint) error {
	avg, err := s.review.GetAvgRating(medicineID)
	if err != nil {
		return err
	}
 
	return s.review.UpdateAvgRating(medicineID, avg)
}

func (s *reviewService) checkMedicineExists(medicineID uint) error {
	if _, err := s.medicine.GetByID(medicineID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMedicineNotFound
		}
		return err
	}
	return nil
}