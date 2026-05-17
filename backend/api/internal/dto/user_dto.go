package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	LastName string `json:"lastname" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}
