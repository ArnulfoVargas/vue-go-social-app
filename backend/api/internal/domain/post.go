package domain

import (
	"Server/internal/dto"
	"Server/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PostService interface {
	CreatePost(userId string, post dto.PostAdd) (model.Post, error)
	GetPost(postId string) (model.Post, error)
	DeletePost(postId string) error
	UpdatePost(postId string, content dto.UpdatePostRequest) (model.Post, error)
	GetPostsByUserId(userId string) ([]model.Post, error)
	ToggleLike(postId string, userId string) error
	GetSuggestedPosts(userId string, limit int) ([]model.Post, error)
}

type PostRepository interface {
	CreatePost(post model.Post) error
	GetPost(postId primitive.ObjectID) (model.Post, error)
	DeletePost(postId primitive.ObjectID) error
	UpdatePost(postId primitive.ObjectID, update bson.M) (model.Post, error)
	GetPostsByUserId(userId primitive.ObjectID) ([]model.Post, error)
	ExistsById(postId primitive.ObjectID) (bool, error)
	GetSuggestedPosts(userId primitive.ObjectID, limit int) ([]model.Post, error)
}
