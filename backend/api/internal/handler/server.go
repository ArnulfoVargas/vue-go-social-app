package handler

import (
	"Server/internal/store"
	"os"

	"github.com/gofiber/fiber/v3"
)

type Server struct {
	app *fiber.App
	db  *store.Database
}

func NewServer(db *store.Database) *Server {
	app := fiber.New()

	return &Server{app: app, db: db}
}

func (s *Server) RegisterRoutes() {
	s.RegisterMainRoutes()
}

func (s *Server) Start() {
	panic(s.app.Listen(":" + os.Getenv("PORT")))
}

func (s *Server) RegisterMainRoutes() {
	s.app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
}
