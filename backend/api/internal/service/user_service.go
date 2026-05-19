package service

import (
	"Server/internal/constants"
	"Server/internal/domain"
	"Server/internal/dto"
	"Server/internal/helpers"
	"Server/internal/model"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type userService struct {
	userRepo   domain.UserRepository
	followRepo domain.FollowRepository
}

func NewUserService(userRepository domain.UserRepository, followRepository domain.FollowRepository) *userService {
	return &userService{
		userRepo:   userRepository,
		followRepo: followRepository,
	}
}

func (s *userService) GetUser(id string) (*model.User, error) {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *userService) UpdateUser(id string, user *dto.UpdateProfileRequest) error {
	data := bson.M{}

	if user.Name != nil {
		data["name"] = *user.Name
	}

	if user.ImageUrl != nil {
		data["image"] = *user.ImageUrl
	}
	if user.Bio != nil {
		data["bio"] = *user.Bio
	}

	return s.userRepo.UpdateUserById(id, data)
}

func (s *userService) GetSuggestedUsers(id string) ([]model.User, error) {
	followingIds, err := s.followRepo.GetFollowingIds(id)
	if err != nil {
		return nil, err
	}

	userId, err := helpers.ToObjectID(id)
	if err != nil {
		return nil, err
	}

	suggestedIds, err := s.followRepo.GetRelatedFollowSuggestions(userId, followingIds, constants.MAX_SUGGESTED_IDS)
	if err != nil {
		return nil, err
	}

	sugIdsLen := len(suggestedIds)
	if sugIdsLen < constants.MAX_SUGGESTED_IDS {
		excludedIds := append(followingIds, userId)
		excludedIds = append(excludedIds, suggestedIds...)

		randomUsers, err := s.userRepo.GetIdsExcluding(excludedIds, constants.MAX_SUGGESTED_IDS-sugIdsLen)

		if err != nil {
			return nil, err
		}

		return s.userRepo.GetUsersByIds(append(suggestedIds, randomUsers...))
	}

	return s.userRepo.GetUsersByIds(suggestedIds)
}

func (s *userService) DeleteUser(id string) error {
	return s.userRepo.DeleteUserById(id)
}
