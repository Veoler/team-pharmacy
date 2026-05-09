package models

import "gorm.io/gorm"

type Status string

const (
	StatusDraft          Status = "draft"
	StatusPendingPayment Status = "pending_payment"
	StatusPaid           Status = "paid"
	StatusCanceled       Status = "canceled"
	StatusShipped        Status = "shipped"
	StatusCompleted      Status = "completed"
)

type Order struct {
	gorm.Model
	User            *User       `json:"-"`
	UserID          uint        `json:"user_id" gorm:"not null;index"`
	Status          Status      `json:"status"`
	Items           []OrderItem `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TotalPrice      int         `json:"total_price"`
	DiscountTotal   int         `json:"discount_total"`
	FinalPrice      int         `json:"final_price"`
	DeliveryAddress string      `json:"delivery_address"`
	Comment         string      `json:"comment"`
}

type OrderCreateRequest struct {
	User            *User       `json:"-"`
	UserID          *uint       `json:"user_id"`
	Status          Status      `json:"status"`
	Items           []OrderItem `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TotalPrice      int         `json:"total_price"`
	DiscountTotal   int         `json:"discount_total"`
	FinalPrice      int         `json:"final_price"`
	DeliveryAddress string      `json:"delivery_address"`
	Comment         string      `json:"comment"`
}

type OrderItem struct {
	gorm.Model
	OrderID      uint   `json:"-" gorm:"not null;index"`
	MedicineID   uint   `json:"medicine_id" gorm:"not null; index"`
	MedicineName string `json:"medicine_name"`
	Quantity     int    `json:"quantity"`
	PricePerUnit int    `json:"price_per_unit"`
	LineTotal    int    `json:"line_total"`
}
