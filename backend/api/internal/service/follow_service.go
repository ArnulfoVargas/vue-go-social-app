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

func (s *followService) ToggleFollowUser(followerID, followingID string) (bool, error) {
	exists, err := s.userRepo.UserExistsById(followingID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, fmt.Errorf("user not found")
	}

	isFollowing, err := s.followRepo.UserIsFollowing(followerID, followingID)

	if err != nil {
		return false, err
	}

	if isFollowing {
		return false, s.followRepo.UnfollowUser(followerID, followingID)
	}

	return true, s.followRepo.FollowUser(followerID, followingID)
}
