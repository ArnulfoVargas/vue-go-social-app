package server

import (
	"Server/internal/repository"
	"Server/internal/service"
)

func (server *Server) RegisterServices() {
	authRepository := repository.NewAuthRepository(server.Db)
	authService := service.NewAuthService(authRepository)
	server.RegisterAuthService(authService)

	userRepository := repository.NewUserRepository(server.Db)
	userService := service.NewUserService(userRepository)
	server.RegisterUserService(userService)
}
