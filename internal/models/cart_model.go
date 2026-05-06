package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	User   *User      `json:"-"`
	UserID uint       `json:"user_id" gorm:"not null;index"`
	Items  []CartItem `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CartCreateUpdateRequest struct {
	User   *User      `json:"-"`
	UserID uint       `json:"user_id"`
	Items  []CartItem `json:"items"`
}

type CartItem struct {
	ID     uint `json:"item_id" gorm:"primaryKey"`
	CartID uint `json:"-" gorm:"not null;index"`
	// Medicine     Medicine `json:"-"`
	MedicineID   uint   `json:"medicine_id" gorm:"not null;index"`
	MedicineName string `json:"medicine_name"`
	Quantity     int    `json:"quantity"`
	PricePerUnit int    `json:"price_per_unit"`
	LineTotal    int    `json:"line_total"`
}
