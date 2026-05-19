package domain

import (
	"Server/internal/dto"
	"Server/internal/model"

	"github.com/gofiber/fiber/v3"
)

type AuthService interface {
	Register(c fiber.Ctx, req dto.RegisterRequest) (string, string, error)
	Login(c fiber.Ctx, req dto.LoginRequest) (string, string, error)
}

type AuthRepository interface {
	CreateUser(req model.User) (*model.User, error)
	FindUserByEmail(email string) (*model.User, error)
}
