package dto

import "github.com/HasanNugroho/gin-clean/pkg/constants"

type (
	CreateUserRequest struct {
		Name        string         `json:"name" validate:"required"`
		Email       string         `json:"email" validate:"required,email"`
		PhoneNumber string         `json:"phone_number" validate:"required"`
		Password    string         `json:"password" validate:"required,min=6"`
		Role        constants.Role `json:"role" validate:"omitempty,oneof=user admin owner customer"`
	}

	UpdateUserRequest struct {
		Name        string         `json:"name,omitempty"`
		Email       string         `json:"email,omitempty" validate:"omitempty,email"`
		PhoneNumber string         `json:"phone_number,omitempty"`
		Password    string         `json:"password,omitempty" validate:"omitempty,min=6"`
		Role        constants.Role `json:"role,omitempty" validate:"omitempty,oneof=user admin owner customer"`
		IsActive    bool           `json:"is_active,omitempty"`
	}

	RegisterUserRequest struct {
		Name        string `json:"name" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number" validate:"required"`
		Password    string `json:"password" validate:"required,min=6"`
	}
)
