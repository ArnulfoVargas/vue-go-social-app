package service

import (
	"Server/internal/domain"
	"fmt"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *userService {
	return &userService{
		repo: userRepository,
	}
}

func (s *userService) GetUser(id string) (*domain.User, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
