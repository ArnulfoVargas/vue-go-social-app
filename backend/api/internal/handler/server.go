package handler

import "github.com/gofiber/fiber/v3"

type Server struct {
	app *fiber.App
}

func NewServer() *Server {
	return &Server{app: fiber.New()}
}

func (s *Server) RegisterRoutes() {
	s.RegisterMainRoutes()
}

func (s *Server) Start() {
	panic(s.app.Listen(":3000"))
}

func (s *Server) RegisterMainRoutes() {
	s.app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
}
