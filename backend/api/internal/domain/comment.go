package domain

import (
	"Server/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CommentService interface {
	AddComment(postId string, userId string, content string) error
	GetComments(postId string) ([]model.Comment, error)
	DeleteComment(commentId string, userId string) error
	UpdateComment(commentId string, userId string, content string) error
}

type CommentRepository interface {
	AddComment(comment model.Comment) error
	GetComments(postId string) ([]model.Comment, error)
	DeleteComment(commentId string) error
	UpdateComment(commentId string, object bson.M) error
}
