package follows

import (
	"Server/internal/helpers"
	"Server/internal/shared"
	"Server/internal/validator"

	"github.com/gofiber/fiber/v3"
)

type FollowHandler struct {
	followService FollowService
}

func NewFollowsHandler(validator *validator.Validator, service FollowService) *FollowHandler {
	return &FollowHandler{
		followService: service,
	}
}

func SetupFollowsRoutes(r fiber.Router, handler *FollowHandler) {
	g := r.Group("/follow", shared.Protected(shared.ParseJWT))

	g.Patch("/:id", handler.toggleFollowUser)
}

// toggleFollowUser toggles the follow status between two users
// @Param id path string true "user ID to follow/unfollow"
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/users/follow/{id} [put]
// @Security BearerAuth
// @Accept json
// @Produce json
// @Tags users
// @Summary Toggle follow status between two users
// @Description Toggles the follow status between two users. If the user is already following the target, unfollows them. If not, follows them.
func (h *FollowHandler) toggleFollowUser(c fiber.Ctx) error {
	id, ok := helpers.GetUserIdFromLocals(c)

	if id == "" || !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(shared.ErrorResponse{
			Error:  "authentication required",
			Status: fiber.StatusUnauthorized,
		})
	}

	pathID := c.Params("id", "")
	if pathID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "user ID is required",
			Status: fiber.StatusBadRequest,
		})
	}

	if id == pathID {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "cannot follow yourself",
			Status: fiber.StatusBadRequest,
		})
	}

	if follow, err := h.followService.ToggleFollowUser(id, pathID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
			Error:  err.Error(),
			Status: fiber.StatusInternalServerError,
		})
	} else {
		message := "follow status toggled to follow"
		if !follow {
			message = "follow status toggled to unfollow"
		}
		return c.JSON(shared.MessageResponse{
			Message: message,
		})
	}
}
