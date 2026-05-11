package models

import (
	"time"

	"gorm.io/gorm"
)
type DisType string

const (
	DisTypePercent	DisType = "percent"	
	DisTypeFixed	DisType = "fixed"	
)


type Promocode struct {
	gorm.Model
	Code			string		`json:"code" gorm:"type:varchar(64);uniqueIndex;not null"`	// строковый код (например, SPRING2025);
	Description		string		`json:"description" gorm:"type:text"`
	DiscountType	DisType		`json:"discount_type"  gorm:"type:varchar(16);not null"`	// percent или fixed;
	DiscountValue	int			`json:"discount_value" gorm:"not null;default:0"`			// размер скидки (процент или фиксированная сумма);
	ValidFrom		time.Time	`json:"valid_from" gorm:"not null"` 
	ValidTo			time.Time	`json:"valid_to" gorm:"not null"` 							// период действия;
	MaxUses			*uint		`json:"max_uses" gorm:"default:0"`							// необязательное ограничение общего количества использований;
	MaxUsesPerUser	*uint		`json:"max_uses_per_user" gorm:"default:0"`					// необязательное ограничение на пользователя;
	IsActive		bool		`json:"is_active" gorm:"not null;default:true"`				// включён/выключен.
}

type PromocodeCreateRequest struct {
    Code            string    `json:"code"`
    Description     string    `json:"description"`
    DiscountType    DisType   `json:"discount_type"`
    DiscountValue   int       `json:"discount_value"`
    ValidFrom       time.Time `json:"valid_from"`
    ValidTo         time.Time `json:"valid_to"`
    MaxUses         *uint     `json:"max_uses,omitempty"`
    MaxUsesPerUser  *uint     `json:"max_uses_per_user,omitempty"`
    IsActive        bool      `json:"is_active"`
}

type PromocodeUpdateRequest struct {
    Code            *string    `json:"code,omitempty"`
    Description     *string    `json:"description,omitempty"`
    DiscountType    *DisType   `json:"discount_type,omitempty"`
    DiscountValue   *int       `json:"discount_value,omitempty"`
    ValidFrom       *time.Time `json:"valid_from,omitempty"`
    ValidTo         *time.Time `json:"valid_to,omitempty"`
    MaxUses         *uint      `json:"max_uses,omitempty"`
    MaxUsesPerUser  *uint      `json:"max_uses_per_user,omitempty"`
    IsActive        *bool      `json:"is_active,omitempty"`
}