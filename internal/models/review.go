package models

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	User		*User		`json:"-"`
	UserID		uint		`json:"user_id" gorm:"not null;index"`
	Medicine	*Medicine	`json:"-"`
	MedicineID 	uint		`json:"medicine_id" gorm:"not null;index"`
	Rating		uint		`json:"rating" gorm:"not null"`				// целое число от 1 до 5;
	Text		string		`json:"text" gorm:"type:text;not null"`		// текст отзыва;
}

type ReviewCreateRequest struct {
	UserID		uint		`json:"user_id"`
	MedicineID 	uint		`json:"-"`
	Rating		uint		`json:"rating"`
	Text		string		`json:"text"`
}

type ReviewUpdateRequest struct {
	Rating		*uint		`json:"rating,omitempty"`
	Text		*string		`json:"text,omitempty"`
}

