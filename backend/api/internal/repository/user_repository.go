package repository

import (
	"Server/internal/domain"
	"Server/internal/store"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userRepository struct {
	db *store.Database
}

func NewUserRepository(db *store.Database) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUser(id string) (*domain.User, error) {
	col := r.db.Database.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user domain.User
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}

	err = col.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)

	if err != nil {
		println(err.Error())
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &user, nil
}
