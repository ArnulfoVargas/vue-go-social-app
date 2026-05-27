package comments

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CommentService interface {
	AddComment(postId, userId, content string) error
	GetComments(postId string) ([]Comment, error)
	DeleteComment(commentId, userId string) error
	UpdateComment(commentId, userId, content string) error
}

type CommentRepository interface {
	AddComment(comment Comment) error
	GetComments(postId primitive.ObjectID) ([]Comment, error)
	DeleteComment(commentId primitive.ObjectID) error
	UpdateComment(commentId primitive.ObjectID, object bson.M) error
}
