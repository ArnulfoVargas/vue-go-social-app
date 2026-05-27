package server

import (
	"Server/internal/features/auth"
	"Server/internal/features/comments"
	"Server/internal/features/follows"
	"Server/internal/features/media"
	domain "Server/internal/features/posts"
	"Server/internal/features/suggestion"
	"Server/internal/features/users"
	"Server/internal/store"
	"Server/internal/validator"
	"os"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
)

type Server struct {
	App               *fiber.App
	Db                *store.Database
	Validator         *validator.Validator
	AuthService       auth.AuthService
	UserService       users.UserService
	FollowService     follows.FollowService
	PostService       domain.PostService
	MediaService      media.MediaService
	SuggestionService suggestion.SuggestionService
	CommentService    comments.CommentService
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
