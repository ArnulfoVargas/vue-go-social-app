package dto

import "Server/internal/model"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type UpdateProfileRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=3,max=50"`
	LastName *string `json:"lastname,omitempty" validate:"omitempty,min=3,max=50"`
	ImageUrl *string `json:"image,omitempty" validate:"omitempty,url"`
	Bio      *string `json:"bio,omitempty" validate:"omitempty,min=3,max=500"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}

type UpdatedProfileResponse struct {
	User    model.User `json:"user"`
	Message string     `json:"message"`
}

type GenericResponse[T any] struct {
	Data T `json:"data"`
}
