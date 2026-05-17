package domain

import (
	"Server/internal/dto"

	"github.com/gofiber/fiber/v3"
)

type AuthService interface {
	Register(c fiber.Ctx, req dto.RegisterRequest) (string, error)
	Login(c fiber.Ctx, req dto.LoginRequest) (string, error)
}

type AuthRepository interface {
	CreateUser(req User) (*User, error)
	FindUserByEmail(email string) (*User, error)
}
