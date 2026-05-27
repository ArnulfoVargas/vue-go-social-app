package server

import (
	"Server/internal/features/auth"
	"Server/internal/features/follows"
	"Server/internal/features/likes"
	"Server/internal/features/media"
	"Server/internal/features/posts"
	"Server/internal/features/suggestion"
	"Server/internal/features/users"
)

func (server *Server) RegisterServices() {
	server.RegisterAuthService()
	server.RegisterUserService()
	server.RegisterFollowService()
	server.RegisterPostService()
	server.RegisterMediaService()
	server.RegisterCommentService()
	server.RegisterSuggestionService()
}

func (server *Server) RegisterUserService() {
	userRepository := users.NewUserRepository(server.Db)
	server.UserService = users.NewUserService(userRepository)
}

func (server *Server) RegisterMediaService() {
	mediaService := media.NewMediaService()
	server.MediaService = mediaService
}

func (server *Server) RegisterAuthService() {
	authRepository := auth.NewAuthRepository(server.Db)
	authService := auth.NewAuthService(authRepository)
	server.AuthService = authService
}

func (server *Server) RegisterFollowService() {
	followRepository := follows.NewFollowRepository(server.Db)
	userRepository := users.NewUserRepository(server.Db)
	followService := follows.NewFollowService(userRepository, followRepository)
	server.FollowService = followService
}

func (server *Server) RegisterPostService() {
	postRepository := posts.NewPostRepository(server.Db)
	likeRepository := likes.NewlikeRepository(server.Db)
	userRepository := users.NewUserRepository(server.Db)
	server.PostService = posts.NewPostService(postRepository, likeRepository, userRepository)
}

func (server *Server) RegisterCommentService() {
}

func (server *Server) RegisterSuggestionService() {
	userRespo := users.NewUserRepository(server.Db)
	followRespo := follows.NewFollowRepository(server.Db)
	suggestionService := suggestion.NewSuggestionService(userRespo, followRespo)
	server.SuggestionService = suggestionService
}
