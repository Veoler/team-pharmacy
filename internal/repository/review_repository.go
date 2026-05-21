package repository

import (
	"github.com/Veoler/team-pharmacy/internal/models"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *models.Review) error
	GetFromMedicine(medicineID uint)([]models.Review, error)
	GetByID(id uint)(*models.Review, error)
	Update(review *models.Review) error
	Delete(id uint) error
	GetAvgRating(medicineID uint) (float64, error)
	UpdateAvgRating(medicineID uint, avg float64) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *models.Review) error {
	if review == nil {
		return nil
	}

	return r.db.Create(review).Error
}

func (r *reviewRepository) GetFromMedicine(medicineID uint)([]models.Review, error) {
	var reviews []models.Review

	if err := r.db.
	Model(&models.Review{}).
	Where("medicine_id = ?", medicineID).
	Find(&reviews).Error; err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *reviewRepository) GetByID(id uint)(*models.Review, error) {
	var review models.Review

	if err := r.db.First(&review, id).Error; err != nil {
		return nil, err
	}

	return &review, nil
}



func (r *reviewRepository) Update(review *models.Review) error {
	if review == nil {
		return nil
	}

	return r.db.Save(review).Error
}

func (r *reviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}

func (r *reviewRepository) GetAvgRating(medicineID uint) (float64, error) {
	var avg float64
 
	if err := r.db.
		Model(&models.Review{}).
		Where("medicine_id = ?", medicineID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avg).Error; err != nil {
		return 0, err
	}
 
	return avg, nil
}

func (r *reviewRepository) UpdateAvgRating(medicineID uint, avg float64) error {
	return r.db.
		Model(&models.Medicine{}).
		Where("id = ?", medicineID).
		Update("avg_rating", avg).Error
}