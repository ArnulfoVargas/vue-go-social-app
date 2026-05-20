package server

import (
	"Server/internal/repository"
	"Server/internal/service"
)

func (server *Server) RegisterServices() {
	authRepository := repository.NewAuthRepository(server.Db)
	authService := service.NewAuthService(authRepository)
	server.RegisterAuthService(authService)

	ufollowRepository := repository.NewFollowRepository(server.Db)
	userRepository := repository.NewUserRepository(server.Db)
	userService := service.NewUserService(userRepository, ufollowRepository)
	server.RegisterUserService(userService)

	followRepository := repository.NewFollowRepository(server.Db)
	fuserRepository := repository.NewUserRepository(server.Db)
	followService := service.NewFollowService(fuserRepository, followRepository)
	server.RegisterFollowService(followService)

	postRepository := repository.NewPostRepository(server.Db)
	likeRepository := repository.NewlikeRepository(server.Db)
	postService := service.NewPostService(postRepository, likeRepository)
	server.RegisterPostService(postService)

	mediaService := service.NewMediaService()
	server.RegisterMediaService(mediaService)
}
