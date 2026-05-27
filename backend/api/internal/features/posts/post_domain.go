package posts

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PostService interface {
	CreatePost(userId string, post PostAdd) (Post, error)
	GetPost(postId string) (Post, error)
	DeletePost(postId string) error
	UpdatePost(postId string, content UpdatePostRequest) (Post, error)
	GetPostsByUserId(userId string) ([]Post, error)
	ToggleLike(postId string, userId string) error
	GetSuggestedPosts(userId string, limit int) ([]Post, error)
}

type PostRepository interface {
	CreatePost(post Post) error
	GetPost(postId primitive.ObjectID) (Post, error)
	DeletePost(postId primitive.ObjectID) error
	UpdatePost(postId primitive.ObjectID, update bson.M) (Post, error)
	GetPostsByUserId(userId primitive.ObjectID) ([]Post, error)
	ExistsById(postId primitive.ObjectID) (bool, error)
	GetSuggestedPosts(userId primitive.ObjectID, limit int) ([]Post, error)
}
