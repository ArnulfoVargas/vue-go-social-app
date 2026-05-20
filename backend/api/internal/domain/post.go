package domain

import (
	"Server/internal/dto"
	"Server/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PostService interface {
	CreatePost(userId string, post dto.PostRequest) (model.Post, error)
	GetPost(postId string) (model.Post, error)
	DeletePost(postId string) error
	UpdatePost(postId string, post dto.PostRequest) (model.Post, error)
	GetPostsByUserId(userId string) ([]model.Post, error)
	AttachImage(postId string, image *model.Media) error
	AttachManyImages(postId string, images []*model.Media) error
	ToggleLike(postId string, userId string) error
	DetachImage(postId string, imageId string) error
	DetachManyImages(postId string, imageIds []string) error
}

type PostRepository interface {
	CreatePost(post model.Post) error
	GetPost(postId string) (model.Post, error)
	DeletePost(postId string) error
	UpdatePost(postId string, update bson.M) (model.Post, error)
	GetPostsByUserId(userId string) ([]model.Post, error)
	AttachImage(postId string, image *model.Media) error
	AttachManyImages(postId string, images []*model.Media) error
	DetachImage(postId string, imageId string) error
	DetachManyImages(postId string, imageIds []string) error
}
