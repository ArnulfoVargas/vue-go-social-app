package server

import (
	"Server/internal/repository"
	"Server/internal/service"
)

func (server *Server) RegisterServices() {
	server.RegisterAuthService()
	server.RegisterUserService()
	server.RegisterFollowService()
	server.RegisterPostService()
	server.RegisterMediaService()
}

func (server *Server) RegisterUserService() {
	followRepository := repository.NewFollowRepository(server.Db)
	userRepository := repository.NewUserRepository(server.Db)
	server.UserService = service.NewUserService(userRepository, followRepository)
}

func (server *Server) RegisterMediaService() {
	mediaService := service.NewMediaService()
	server.MediaService = mediaService
}

func (server *Server) RegisterAuthService() {
	authRepository := repository.NewAuthRepository(server.Db)
	authService := service.NewAuthService(authRepository)
	server.AuthService = authService
}

func (server *Server) RegisterFollowService() {
	followRepository := repository.NewFollowRepository(server.Db)
	fuserRepository := repository.NewUserRepository(server.Db)
	followService := service.NewFollowService(fuserRepository, followRepository)
	server.FollowService = followService
}

func (server *Server) RegisterPostService() {
	postRepository := repository.NewPostRepository(server.Db)
	likeRepository := repository.NewlikeRepository(server.Db)
	userRepository := repository.NewUserRepository(server.Db)
	server.PostService = service.NewPostService(postRepository, likeRepository, userRepository)
}
