package handler

import (
	"Server/internal/domain"
	"Server/internal/dto"
	"Server/internal/validator"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	validator *validator.Validator
	service   domain.AuthService
}

func NewAuthHandler(validator *validator.Validator, service domain.AuthService) *AuthHandler {
	return &AuthHandler{
		validator: validator,
		service:   service,
	}
}

func SetupAuthRoutes(s fiber.Router, validator *validator.Validator, service domain.AuthService) {
	g := s.Group("/auth")

	handler := NewAuthHandler(validator, service)

	g.Post("/register", handler.register)
	g.Post("/login", handler.login)
}

// Register a new user
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration details"
// @Router /api/v1/auth/register [post]
// @Success 200 {object} domain.AuthResponse
func (h *AuthHandler) register(c fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	token, err := h.service.Register(c, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(domain.AuthResponse{Token: token})
}

// Login a user
// @Summary Login a user
// @Description Login a user with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login credentials"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) login(c fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	token, err := h.service.Login(c, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(domain.AuthResponse{Token: token})
}
