package handler

import (
	"Server/internal/constants"
	"Server/internal/server"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(s *server.Server) {
	s.App.Get("/", HelloWorld)

	g := s.App.Group("/api/v1")

	authHandler := NewAuthHandler(s.Validator, s.AuthService)
	SetupAuthRoutes(g, authHandler)

	userHandler := NewUserHandler(s.Validator, s.UserService, s.FollowService)
	SetupUserRoutes(g, userHandler)
	postHandler := NewPostHandler(s.Validator, s.PostService, s.MediaService)
	SetupPostRoutes(g, postHandler)
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

func getUserIdFromLocals(c fiber.Ctx) (string, bool) {
	id, ok := c.Locals(constants.USER_ID_CLAIM).(string)
	return id, ok
}
