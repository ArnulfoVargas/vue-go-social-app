package service

import (
	"Server/internal/domain"
	"Server/internal/helpers"
	"Server/internal/model"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	followerId, err := helpers.ToObjectID(followerID)
	if err != nil {
		return false, err
	}

	followingId, err := helpers.ToObjectID(followingID)
	if err != nil {
		return false, err
	}

	exists, err := s.userRepo.UserExistsById(followerId)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, fmt.Errorf("user not found")
	}

	isFollowing, err := s.followRepo.UserIsFollowing(followerId, followingId)

	if err != nil {
		return false, err
	}

	if isFollowing {
		return false, s.followRepo.UnfollowUser(followerId, followingId)
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	follow := model.Follow{
		ID:          primitive.NewObjectID(),
		FollowerID:  followerId,
		FollowingID: followingId,
		CreatedAt:   now,
		UpdatedAt:   now,
		Status:      1,
	}

	return true, s.followRepo.FollowUser(follow)
}
