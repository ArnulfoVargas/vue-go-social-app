package auth

import (
	"Server/internal/features/users"
	"Server/internal/helpers"
	"Server/internal/store"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type authRepository struct {
	collection *mongo.Collection
}

func NewAuthRepository(db *store.Database) *authRepository {
	return &authRepository{
		collection: db.Database.Collection("users"),
	}
}

func (r *authRepository) CreateUser(req users.User) (*users.User, error) {
	col := r.collection

	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	_, err := col.InsertOne(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &req, nil
}

func (r *authRepository) FindUserByEmail(email string) (*users.User, error) {
	col := r.collection

	var user users.User

	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	err := col.FindOne(ctx, bson.M{"email": email, "status": 1}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	return &user, nil
}
