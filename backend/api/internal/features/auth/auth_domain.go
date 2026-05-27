package auth

import (
	"Server/internal/features/users"

	"github.com/gofiber/fiber/v3"
)

type AuthService interface {
	Register(c fiber.Ctx, req RegisterRequest) (string, string, error)
	Login(c fiber.Ctx, req LoginRequest) (string, string, error)
}

type AuthRepository interface {
	CreateUser(req users.User) (*users.User, error)
	FindUserByEmail(email string) (*users.User, error)
}
