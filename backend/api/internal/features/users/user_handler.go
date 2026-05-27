package users

import (
	"Server/internal/features/media"
	"Server/internal/helpers"
	"Server/internal/shared"
	"Server/internal/validator"
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	validator    *validator.Validator
	userService  UserService
	mediaService media.MediaService
}

func NewUserHandler(validator *validator.Validator, service UserService, mediaService media.MediaService) *UserHandler {
	return &UserHandler{
		mediaService: mediaService,
		validator:    validator,
		userService:  service,
	}
}

func SetupUserRoutes(r fiber.Router, handler *UserHandler) {
	g := r.Group("/users", shared.Protected(shared.ParseJWT))

	g.Post("/profile", handler.AddProfilePicture)
	g.Delete("/profile", handler.RemoveProfilePicture)
	g.Get("/:id", handler.GetUser)
	g.Patch("/:id", handler.UpdateUser)
	g.Delete("/:id", handler.DeleteUser)
}

// @Summary Get user by ID
// @Description Get user by ID
// @Tags users
// @Param id path string true "User ID"
// @Router /api/v1/users/{id} [get]
// @Security BearerAuth
// @Produce json
// @Failure 400 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
func (h *UserHandler) GetUser(c fiber.Ctx) error {
	id := c.Params("id", "")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "id is required",
			Status: fiber.StatusBadRequest,
		})
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(404).JSON(shared.ErrorResponse{
				Error:  "profile not found",
				Status: 404,
			})
		}

		return c.Status(fiber.ErrInternalServerError.Code).JSON(shared.ErrorResponse{
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
// @Param user body UpdateProfileRequest true "User data"
// @Router /api/v1/users/{id} [patch]
// @Security BearerAuth
// @Accept json
// @Produce json
// @Error 400 {object} shared.ErrorResponse "invalid request body"
// @Error 401 {object} shared.ErrorResponse "authentication required"
// @Error 403 {object} shared.ErrorResponse "not authorized to update this profile"
// @Error 500 {object} shared.ErrorResponse "internal server error"
func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id, ok := helpers.GetUserIdFromLocals(c)

	if id == "" || !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(shared.ErrorResponse{
			Error:  "authentication required",
			Status: fiber.StatusUnauthorized,
		})
	}

	pathID := c.Params("id", "")
	if pathID != id {
		return c.Status(fiber.StatusUnauthorized).JSON(shared.ErrorResponse{
			Error:  "not authorized to update this profile",
			Status: fiber.StatusUnauthorized,
		})
	}

	var user UpdateProfileRequest
	if err := c.Bind().Body(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "invalid request body",
			Status: fiber.StatusBadRequest,
		})
	}

	if err := h.validator.Validate(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := h.userService.UpdateUser(id, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
			Error:  err.Error(),
			Status: fiber.StatusInternalServerError,
		})
	}

	updatedUser, err := h.userService.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
			Error:  err.Error(),
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.JSON(UpdatedProfileResponse{
		User:    *updatedUser,
		Message: "profile updated successfully",
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

	localsId, ok := helpers.GetUserIdFromLocals(c)

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

// @Tags users
// @Router /api/v1/users/{id}/profile-picture [post]
// @Security BearerAuth
// @Summary Add profile picture
// @Description Adds a profile picture to the user with the given ID
// @Param file formData file true "Profile picture"
// @Success 200 {object} SetProfilePictureResponse
func (h *UserHandler) AddProfilePicture(c fiber.Ctx) error {
	id, ok := helpers.GetUserIdFromLocals(c)
	if id == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var files []*multipart.FileHeader
	form, err := c.MultipartForm()
	if err == nil {
		files = form.File["images"]
	}

	if len(files) != 1 {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "only one image is allowed",
			Status: fiber.StatusBadRequest,
		})
	}

	var mediafile *media.Media
	if len(files) == 1 {
		mediafile, err = h.mediaService.Upload(files[0])
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
				Error:  "failed to upload images",
				Status: fiber.StatusInternalServerError,
			})
		}
	}

	if mediafile == nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "no image provided",
			Status: fiber.StatusBadRequest,
		})
	}

	if err := h.userService.AddProfilePicture(id, *mediafile); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(SetProfilePictureResponse{
		AvatarUrl: media.MediaResponse{
			ID:  mediafile.ID.Hex(),
			URL: mediafile.URL,
		},
	})
}

// @Tags users
// @Router /api/v1/users/{id}/profile-picture [delete]
// @Security BearerAuth
// @Summary Remove profile picture
// @Description Removes the profile picture from the user with the given ID
// @Param id path string true "User ID"
func (h *UserHandler) RemoveProfilePicture(c fiber.Ctx) error {
	id, ok := helpers.GetUserIdFromLocals(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if err := h.userService.RemoveProfilePicture(id); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
