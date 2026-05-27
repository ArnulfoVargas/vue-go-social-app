package users

import (
	"Server/internal/features/media"
	"Server/internal/helpers"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepository UserRepository) *userService {
	return &userService{
		userRepo: userRepository,
	}
}

func (s *userService) GetUser(id string) (*User, error) {
	uId, err := helpers.ToObjectID(id)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetUserById(uId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *userService) UpdateUser(id string, user *UpdateProfileRequest) error {
	uId, err := helpers.ToObjectID(id)
	if err != nil {
		return err
	}
	data := bson.M{}

	if user.Name != nil {
		data["name"] = *user.Name
	}

	if user.Bio != nil {
		data["bio"] = *user.Bio
	}

	return s.userRepo.UpdateUserById(uId, data)
}

func (s *userService) DeleteUser(id string) error {
	uId, err := helpers.ToObjectID(id)
	if err != nil {
		return err
	}
	return s.userRepo.DeleteUserById(uId)
}

func (s *userService) AddProfilePicture(id string, media media.Media) error {
	uId, err := helpers.ToObjectID(id)
	if err != nil {
		return err
	}

	return s.userRepo.SetProfilePicture(uId, media)
}

func (s *userService) RemoveProfilePicture(id string) error {
	uId, err := helpers.ToObjectID(id)
	if err != nil {
		return err
	}
	return s.userRepo.RemoveProfilePicture(uId)
}

func (s *userService) ExistsUser(id string) (bool, error) {
	uId, err := helpers.ToObjectID(id)
	if err != nil {
		return false, err
	}
	return s.userRepo.UserExistsById(uId)
}
