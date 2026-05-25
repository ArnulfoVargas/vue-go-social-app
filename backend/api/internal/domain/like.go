package domain

import (
	"Server/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeRepository interface {
	DeleteLike(like model.Like) error
	AddLike(like model.Like) error
	HasLike(postId, userId primitive.ObjectID) (bool, error)
	DeleteLikesFromPost(postId primitive.ObjectID) error
}
