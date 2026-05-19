package handler

import (
	"Server/internal/server"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(s *server.Server) {
	s.App.Get("/", HelloWorld)

	g := s.App.Group("/api/v1")
	SetupAuthRoutes(g, s.Validator, s.AuthService)
	SetupUserRoutes(g, s.Validator, s.UserService, s.FollowService)
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
