package domain

import (
	"Server/internal/dto"
	"Server/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserService interface {
	GetUser(id string) (*model.User, error)
	UpdateUser(id string, user *dto.UpdateProfileRequest) error
	GetSuggestedUsers(id string) ([]model.User, error)
}

type UserRepository interface {
	GetUserById(id string) (*model.User, error)
	UpdateUserById(id string, data bson.M) error
	UserExistsById(id string) (bool, error)
	GetUsersExcluding(excludeIDs []primitive.ObjectID, limit int) ([]model.User, error)
	GetUsersByIds(ids []primitive.ObjectID) ([]model.User, error)
	GetIdsExcluding(excludeIDs []primitive.ObjectID, limit int) ([]primitive.ObjectID, error)
}
