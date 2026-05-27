package auth

import (
	"Server/internal/shared"
	"Server/internal/validator"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	validator *validator.Validator
	service   AuthService
}

func NewAuthHandler(validator *validator.Validator, service AuthService) *AuthHandler {
	return &AuthHandler{
		validator: validator,
		service:   service,
	}
}

func SetupAuthRoutes(s fiber.Router, handler *AuthHandler) {
	g := s.Group("/auth")

	g.Post("/register", handler.register)
	g.Post("/login", handler.login)
}

// Register a new user
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration details"
// @Router /api/v1/auth/register [post]
// @Success 200 {object} AuthResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
func (h *AuthHandler) register(c fiber.Ctx) error {
	var req RegisterRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "invalid body",
			Status: fiber.StatusBadRequest,
		})
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	token, id, err := h.service.Register(c, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
			Error:  err.Error(),
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.JSON(AuthResponse{Token: token, Id: id})
}

// Login a user
// @Summary Login a user
// @Description Login a user with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login credentials"
// @Router /api/v1/auth/login [post]
// @Success 200 {object} AuthResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
func (h *AuthHandler) login(c fiber.Ctx) error {
	var req LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "invalid body",
			Status: fiber.StatusBadRequest,
		})
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	token, id, err := h.service.Login(c, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
			Error:  err.Error(),
			Status: fiber.ErrInternalServerError.Code,
		})
	}

	return c.JSON(AuthResponse{Token: token, Id: id})
}
