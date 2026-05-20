package repository

import (
	"Server/internal/store"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type likeRepository struct {
	collection *mongo.Collection
}

func NewlikeRepository(db *store.Database) *likeRepository {
	return &likeRepository{collection: db.Database.Collection("likes")}
}

func (r *likeRepository) ToggleLike(postId string, userId string) error {
	return nil
}
