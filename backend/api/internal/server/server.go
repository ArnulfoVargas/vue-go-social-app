package server

import (
	"Server/internal/domain"
	"Server/internal/store"
	"Server/internal/validator"
	"os"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
)

type Server struct {
	App           *fiber.App
	Db            *store.Database
	Validator     *validator.Validator
	AuthService   domain.AuthService
	UserService   domain.UserService
	FollowService domain.FollowService
	PostService   domain.PostService
	MediaService  domain.MediaService
}

var (
	swaggerConfig = swaggo.Config{
		Title: "Social Media API",
	}
)

func (s *Server) UseSwagger() {
	s.App.Get("/swagger/*", swaggo.New(swaggerConfig))
}

func NewServer(db *store.Database) *Server {
	app := fiber.New()

	server := Server{App: app, Db: db, Validator: validator.New()}

	return &server
}

func (s *Server) Start() {
	panic(s.App.Listen(":" + os.Getenv("PORT")))
}
