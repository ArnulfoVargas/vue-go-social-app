package handler

import (
	"Server/internal/store"
	"os"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
)

type Server struct {
	app *fiber.App
	db  *store.Database
}

var (
	swaggerConfig = swaggo.Config{
		Title: "Social Media API",
	}
)

func (s *Server) UseSwagger() {
	s.app.Get("/swagger/*", swaggo.New(swaggerConfig))
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
	s.app.Get("/", s.HelloWorld)
}

// @Summary Hello World
// @Description Returns a simple greeting
// @Tags Testing
// @Produce plain
// @Success 200 {string} string "Hello, World!"
// @Router / [get]
func (s *Server) HelloWorld(c fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
