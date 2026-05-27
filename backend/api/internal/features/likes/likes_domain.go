package likes

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeRepository interface {
	DeleteLike(like Like) error
	AddLike(like Like) error
	HasLike(postId, userId primitive.ObjectID) (bool, error)
	DeleteLikesFromPost(postId primitive.ObjectID) error
}
