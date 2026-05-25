package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `json:"name"`
}

type CategoryCreateRequest struct {
	Name *string `json:"name"`
}

////////////////////////////////////////////////////////////
type Subcategory struct {
	gorm.Model
	CategoryID uint   `json:"category_id" gorm:"not null;index"`
	Name       string `json:"name"`
}

type SubcategoryCreateRequest struct {
	CategoryID uint   `json:"subcategory_id"`
	Name       string `json:"name"`
}
