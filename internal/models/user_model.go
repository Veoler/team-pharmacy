package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName       string  `json:"full_name"`
	Email          string  `json:"e-mail"`
	Phone          string  `json:"phone"`
	DefaultAddress string  `json:"default_address"`
	Orders         []Order `json:"-"`
}

type UserCreateRequest struct {
	FullName       *string `json:"full_name"`
	Email          *string `json:"e-mail"`
	Phone          string  `json:"phone"`
	DefaultAddress string  `json:"default_address"`
}

type UserUpdateRequest struct {
	FullName       *string `json:"full_name"`
	Email          *string `json:"e-mail"`
	Phone          *string `json:"phone"`
	DefaultAddress *string `json:"default_address"`
}
