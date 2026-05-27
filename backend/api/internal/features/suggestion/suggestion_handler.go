package suggestion

import (
	"Server/internal/features/users"
	"Server/internal/helpers"
	"Server/internal/shared"
	"Server/internal/validator"

	"github.com/gofiber/fiber/v3"
)

type SuggestionHandler struct {
	suggestionService SuggestionService
}

func NewSuggestionHandler(validator *validator.Validator, service SuggestionService) *SuggestionHandler {
	return &SuggestionHandler{
		suggestionService: service,
	}
}

func SetupSuggestionRoutes(r fiber.Router, handler *SuggestionHandler) {
	g := r.Group("/users", shared.Protected(shared.ParseJWT))

	g.Get("/suggest", handler.GetSuggestedUsers)
}

// @Summary Get suggested users
// @Description Returns a list of suggested users for the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/users/suggest [get]
// @Security BearerAuth
func (s *SuggestionHandler) GetSuggestedUsers(c fiber.Ctx) error {
	id, ok := helpers.GetUserIdFromLocals(c)
	if id == "" || !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(shared.ErrorResponse{
			Error:  "authentication required",
			Status: fiber.StatusUnauthorized,
		})
	}

	sugUsers, err := s.suggestionService.GetSuggestedUsers(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
			Error:  err.Error(),
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(shared.GenericResponse[[]users.User]{
		Data: sugUsers,
	})
}
