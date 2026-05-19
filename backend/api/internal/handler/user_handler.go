package handler

import (
	"Server/internal/domain"
	"Server/internal/dto"
	"Server/internal/middleware"
	"Server/internal/model"
	"Server/internal/service"
	"Server/internal/validator"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	validator     *validator.Validator
	userService   domain.UserService
	followService domain.FollowService
}

func NewUserHandler(validator *validator.Validator, service domain.UserService, followService domain.FollowService) *UserHandler {
	return &UserHandler{
		validator:     validator,
		userService:   service,
		followService: followService,
	}
}

func SetupUserRoutes(r fiber.Router, v *validator.Validator, s domain.UserService, fs domain.FollowService) {
	g := r.Group("/users", middleware.Protected(service.ParseJWT))

	handler := NewUserHandler(v, s, fs)

	g.Get("/suggest", handler.GetSuggestedUsers)
	g.Get("/:id", handler.GetUser)
	g.Patch("/:id", handler.UpdateUser)
	g.Put("/follow/:id", handler.ToggleFollowUser)
	g.Delete("/:id", handler.DeleteUser)
}

// @Summary Get user by ID
// @Description Get user by ID
// @Tags users
// @Param id path string true "User ID"
// @Router /api/v1/users/{id} [get]
// @Security BearerAuth
// @Produce json
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
func (h *UserHandler) GetUser(c fiber.Ctx) error {
	id := c.Params("id", "")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "id is required",
			Status: fiber.StatusBadRequest,
		})
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(404).JSON(dto.ErrorResponse{
				Error:  "profile not found",
				Status: 404,
			})
		}

		return c.Status(fiber.ErrInternalServerError.Code).JSON(dto.ErrorResponse{
			Error:  err.Error(),
			Status: fiber.ErrInternalServerError.Code,
		})
	}

	return c.JSON(user)
}

// @Summary Update user
// @Description Update user
// @Tags users
// @Param id path string true "User ID"
// @Param user body dto.UpdateProfileRequest true "User data"
// @Router /api/v1/users/{id} [patch]
// @Security BearerAuth
// @Accept json
// @Produce json
// @Error 400 {object} dto.ErrorResponse "invalid request body"
// @Error 401 {object} dto.ErrorResponse "authentication required"
// @Error 403 {object} dto.ErrorResponse "not authorized to update this profile"
// @Error 500 {object} dto.ErrorResponse "internal server error"
func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id, ok := c.Locals("userID").(string)

	if id == "" || !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error:  "authentication required",
			Status: fiber.StatusUnauthorized,
		})
	}

	pathID := c.Params("id", "")
	if pathID != id {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error:  "not authorized to update this profile",
			Status: fiber.StatusUnauthorized,
		})
	}

	var user dto.UpdateProfileRequest
	if err := c.Bind().Body(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "invalid request body",
			Status: fiber.StatusBadRequest,
		})
	}

	if err := h.validator.Validate(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := h.userService.UpdateUser(id, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:  err.Error(),
			Status: fiber.StatusInternalServerError,
		})
	}

	updatedUser, err := h.userService.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:  err.Error(),
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.JSON(dto.UpdatedProfileResponse{
		User:    *updatedUser,
		Message: "profile updated successfully",
	})
}

// ToggleFollowUser toggles the follow status between two users
// @Param id path string true "user ID to follow/unfollow"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/users/follow/{id} [put]
// @Security BearerAuth
// @Accept json
// @Produce json
// @Tags users
// @Summary Toggle follow status between two users
// @Description Toggles the follow status between two users. If the user is already following the target, unfollows them. If not, follows them.
func (h *UserHandler) ToggleFollowUser(c fiber.Ctx) error {
	id, ok := c.Locals("userID").(string)
	if id == "" || !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error:  "authentication required",
			Status: fiber.StatusUnauthorized,
		})
	}

	pathID := c.Params("id", "")
	if pathID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "user ID is required",
			Status: fiber.StatusBadRequest,
		})
	}

	if id == pathID {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "cannot follow yourself",
			Status: fiber.StatusBadRequest,
		})
	}

	if follow, err := h.followService.ToggleFollowUser(id, pathID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:  err.Error(),
			Status: fiber.StatusInternalServerError,
		})
	} else {
		message := "follow status toggled to follow"
		if !follow {
			message = "follow status toggled to unfollow"
		}
		return c.JSON(dto.MessageResponse{
			Message: message,
		})
	}
}

// @Summary Get suggested users
// @Description Returns a list of suggested users for the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Router /api/v1/users/suggest [get]
// @Security BearerAuth
func (h *UserHandler) GetSuggestedUsers(c fiber.Ctx) error {
	id, ok := c.Locals("userID").(string)
	if id == "" || !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error:  "authentication required",
			Status: fiber.StatusUnauthorized,
		})
	}

	users, err := h.userService.GetSuggestedUsers(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:  err.Error(),
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.GenericResponse[[]model.User]{
		Data: users,
	})
}

// @Tags users
// @Router /api/v1/users/{id} [delete]
// @Security BearerAuth
// @Summary Delete user
// @Description Deletes the user with the given ID
// @Param id path string true "User ID"
func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	localsId, ok := c.Locals("userID").(string)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if id != localsId {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err := h.userService.DeleteUser(id); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
