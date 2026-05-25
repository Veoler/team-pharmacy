package models

import "gorm.io/gorm"

type Medicine struct {
	gorm.Model
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Price                int    `json:"price"`
	InStock              bool   `json:"in_stock"`
	StockQuantity        int    `json:"stock_quantity"`
	CategoryID           uint   `json:"category_id" gorm:"not null;index"`
	SubCategoryID        uint   `json:"subcategory_id" gorm:"not null;index"`
	Manufacturer         string `json:"manufacturer"`
	PrescriptionRequired bool   `json:"prescription_required"`
	AvgRating            float64    `json:"avg_rating"`
}

type MedicineInfo struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	InStock   bool   `json:"in_stock"`
	AvgRating float64    `json:"avg_rating"`
}

type MedicineCreateRequest struct {
	Name                 *string `json:"name"`
	Description          *string `json:"description"`
	Price                *int    `json:"price"`
	InStock              *bool   `json:"in_stock"`
	StockQuantity        *int    `json:"stock_quantity"`
	CategoryID           *uint   `json:"category_id"`
	SubCategoryID        *uint   `json:"subcategory_id"`
	Manufacturer         *string `json:"manufacturer"`
	PrescriptionRequired bool    `json:"prescription_required"`
}

type MedicineUpdateRequest struct {
	Name                 *string `json:"name"`
	Description          *string `json:"description"`
	Price                *int    `json:"price"`
	InStock              *bool   `json:"in_stock"`
	StockQuantity        *int    `json:"stock_quantity"`
	CategoryID           *uint   `json:"category_id"`
	SubCategoryID        *uint   `json:"subcategory_id"`
	Manufacturer         *string `json:"manufacturer"`
	PrescriptionRequired *bool   `json:"prescription_required"`
}
