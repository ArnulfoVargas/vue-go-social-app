package repository

import (
	"Server/internal/domain"
	"Server/internal/store"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type authRepository struct {
	db *store.Database
}

func NewAuthRepository(db *store.Database) *authRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(req domain.User) (*domain.User, error) {
	collection := r.db.Database.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &req, nil
}

func (r *authRepository) FindUserByEmail(email string) (*domain.User, error) {
	collection := r.db.Database.Collection("users")
	var user domain.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"email": email, "status": 1}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	return &user, nil
}
