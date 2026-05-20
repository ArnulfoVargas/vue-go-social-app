package repository

import (
	"Server/internal/model"
	"Server/internal/store"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type postRepository struct {
	postCollection      *mongo.Collection
	postMediaCollection *mongo.Collection
}

func NewPostRepository(db *store.Database) *postRepository {
	return &postRepository{
		postCollection:      db.Database.Collection("posts"),
		postMediaCollection: db.Database.Collection("post_media"),
	}
}

func (p *postRepository) CreatePost(post model.Post) error {
	return nil
}

func (p *postRepository) GetPost(postId string) (model.Post, error) {
	return model.Post{}, nil
}

func (p *postRepository) DeletePost(postId string) error {
	return nil
}

func (p *postRepository) UpdatePost(postId string, update bson.M) (model.Post, error) {
	return model.Post{}, nil
}

func (p *postRepository) GetPostsByUserId(userId string) ([]model.Post, error) {
	return nil, nil
}

func (p *postRepository) AttachImage(postId string, image *model.Media) error {
	return nil
}

func (p *postRepository) AttachManyImages(postId string, images []*model.Media) error {
	return nil
}

func (p *postRepository) DetachImage(postId string, imageId string) error {
	return nil
}

func (p *postRepository) DetachManyImages(postId string, imageIds []string) error {
	return nil
}
