package follows

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FollowService interface {
	ToggleFollowUser(userID, targetUserID string) (bool, error)
}

type FollowRepository interface {
	FollowUser(follow Follow) error
	UnfollowUser(userID, targetUserID primitive.ObjectID) error
	UserIsFollowing(userID, targetUserID primitive.ObjectID) (bool, error)
	GetFollowingCount(userID primitive.ObjectID) (int64, error)
	GetFollowerCount(userID primitive.ObjectID) (int64, error)
	GetFollowingIds(userID primitive.ObjectID) ([]primitive.ObjectID, error)
	GetRelatedFollowSuggestions(userId primitive.ObjectID, followingIds []primitive.ObjectID, limit int) ([]primitive.ObjectID, error)
}
