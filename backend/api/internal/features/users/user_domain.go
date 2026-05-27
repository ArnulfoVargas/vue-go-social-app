package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserService interface {
	GetUser(id string) (*User, error)
	UpdateUser(id string, user *UpdateProfileRequest) error
	DeleteUser(id string) error
}

type UserRepository interface {
	GetUserById(id primitive.ObjectID) (*User, error)
	UpdateUserById(id primitive.ObjectID, data bson.M) error
	UserExistsById(id primitive.ObjectID) (bool, error)
	GetUsersExcluding(excludeIDs []primitive.ObjectID, limit int) ([]User, error)
	GetUsersByIds(ids []primitive.ObjectID) ([]User, error)
	GetIdsExcluding(excludeIDs []primitive.ObjectID, limit int) ([]primitive.ObjectID, error)
	DeleteUserById(id primitive.ObjectID) error
}
