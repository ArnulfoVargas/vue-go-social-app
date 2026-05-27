package router

import (
	"Server/internal/features/auth"
	"Server/internal/features/follows"
	"Server/internal/features/posts"
	"Server/internal/features/profile"
	"Server/internal/features/suggestion"
	"Server/internal/features/users"
	"Server/internal/server"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(s *server.Server) {
	s.App.Get("/", HelloWorld)

	g := s.App.Group("/api/v1")

	authHandler := auth.NewAuthHandler(s.Validator, s.AuthService)
	auth.SetupAuthRoutes(g, authHandler)

	userHandler := users.NewUserHandler(s.Validator, s.UserService, s.MediaService)
	users.SetupUserRoutes(g, userHandler)

	postHandler := posts.NewPostHandler(s.Validator, s.PostService, s.MediaService)
	posts.SetupPostRoutes(g, postHandler)

	followsHandler := follows.NewFollowsHandler(s.Validator, s.FollowService)
	follows.SetupFollowsRoutes(g, followsHandler)

	suggestionHandler := suggestion.NewSuggestionHandler(s.Validator, s.SuggestionService)
	suggestion.SetupSuggestionRoutes(g, suggestionHandler)

	profileHandler := profile.NewProfileHandler(s.ProfileService, s.UserService)
	profile.SetupProfileRoutes(g, profileHandler)
}

// @Summary Hello World
// @Description Returns a simple greeting
// @Tags Testing
// @Produce plain
// @Success 200 {string} string "Hello, World!"
// @Router / [get]
func HelloWorld(c fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
