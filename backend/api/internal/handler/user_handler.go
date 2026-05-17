package handler

import (
	"Server/internal/domain"
	"Server/internal/middleware"
	"Server/internal/service"
	"Server/internal/validator"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	validator *validator.Validator
	service   domain.UserService
}

func NewUserHandler(validator *validator.Validator, service domain.UserService) *UserHandler {
	return &UserHandler{
		validator: validator,
		service:   service,
	}
}

func SetupUserRoutes(r fiber.Router, v *validator.Validator, s domain.UserService) {
	g := r.Group("/users", middleware.Protected(service.ParseJWT))

	handler := NewUserHandler(v, s)

	g.Get("/:id", handler.GetUser)
}

// @Summary Get user by ID
// @Description Get user by ID
// @Tags users
// @Param id path string true "User ID"
// @Router /api/v1/users/{id} [get]
// @Security BearerAuth
func (h *UserHandler) GetUser(c fiber.Ctx) error {
	id := c.Params("id", "")

	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "id is required")
	}

	user, err := h.service.GetUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(user)
}
