package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FollowService interface {
	ToggleFollowUser(userID, targetUserID string) (bool, error)
}

type FollowRepository interface {
	FollowUser(userID, targetUserID string) error
	UnfollowUser(userID, targetUserID string) error
	UserIsFollowing(userID, targetUserID string) (bool, error)
	GetFollowingCount(userID string) (int64, error)
	GetFollowerCount(userID string) (int64, error)
	GetFollowingIds(userID string) ([]primitive.ObjectID, error)
	GetRelatedFollowSuggestions(userId primitive.ObjectID, followingIds []primitive.ObjectID, limit int) ([]primitive.ObjectID, error)
}
