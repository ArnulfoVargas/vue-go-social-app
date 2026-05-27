package profile

import (
	"Server/internal/features/users"
	"Server/internal/shared"

	"github.com/gofiber/fiber/v3"
)

type ProfileHandler struct {
	profileService ProfileService
	userService    users.UserService
}

func NewProfileHandler(profileService ProfileService, userService users.UserService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

func SetupProfileRoutes(r fiber.Router, handler *ProfileHandler) {
	g := r.Group("/profile", shared.Protected(shared.ParseJWT))

	g.Get("/:id", handler.GetProfile)
}

// GetProfile returns the profile of a user by their ID
// @Summary Get profile
// @Description Get the profile of a user by their ID
// @Param id path string true "User ID"
// @Produce json
// @Success 200 {object} Profile
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
func (h *ProfileHandler) GetProfile(c fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if exists, err := h.userService.ExistsUser(userID); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	} else if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}

	profile, err := h.profileService.GetProfile(userID)
	if err != nil {
		return err
	}
	return c.JSON(*profile)
}
