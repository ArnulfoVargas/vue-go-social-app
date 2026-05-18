package service

import (
	"Server/internal/domain"
	"fmt"
)

type followService struct {
	userRepo   domain.UserRepository
	followRepo domain.FollowRepository
}

func NewFollowService(userRepo domain.UserRepository, followRepo domain.FollowRepository) *followService {
	return &followService{
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

func (s *followService) ToggleFollowUser(followerID, followingID string) error {
	exists, err := s.userRepo.UserExistsById(followingID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	isFollowing, err := s.followRepo.UserIsFollowing(followerID, followingID)

	if err != nil {
		return err
	}

	if isFollowing {
		return s.followRepo.UnfollowUser(followerID, followingID)
	}

	return s.followRepo.FollowUser(followerID, followingID)
}
