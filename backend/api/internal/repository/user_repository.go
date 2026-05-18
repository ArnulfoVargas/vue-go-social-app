package repository

import (
	"Server/internal/helpers"
	"Server/internal/model"
	"Server/internal/store"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *store.Database) *userRepository {
	return &userRepository{
		collection: db.Database.Collection("users"),
	}
}

func (r *userRepository) GetUserById(userId string) (*model.User, error) {
	col := r.collection

	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	var user model.User
	objId, err := helpers.ToObjectID(userId)
	if err != nil {
		return nil, err
	}

	err = col.FindOne(ctx, bson.M{"_id": objId, "status": 1}).Decode(&user)

	if err != nil {
		println(err.Error())
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) UpdateUserById(userId string, data bson.M) error {
	col := r.collection

	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	objId, err := helpers.ToObjectID(userId)
	if err != nil {
		return err
	}

	data["updatedAt"] = primitive.NewDateTimeFromTime(time.Now())

	_, err = col.UpdateOne(ctx, bson.M{"_id": objId, "status": 1}, bson.M{"$set": data})
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (r *userRepository) UserExistsById(userId string) (bool, error) {
	id, err := helpers.ToObjectID(userId)
	if err != nil {
		return false, err
	}

	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, fmt.Errorf("error finding user")
	}

	return count > 0, nil
}

func (r *userRepository) GetUsersExcluding(excludeIDs []primitive.ObjectID, limit int) ([]model.User, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	filter := bson.M{"status": 1}
	if len(excludeIDs) > 0 {
		filter["_id"] = bson.M{"$nin": excludeIDs}
	}

	pipeline := mongo.Pipeline{
		{{
			Key:   "$match",
			Value: filter,
		}},
		{{
			Key:   "$sample",
			Value: bson.M{"size": limit},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error finding users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("error decoding users")
	}

	return users, nil
}

func (r *userRepository) GetUsersByIds(ids []primitive.ObjectID) ([]model.User, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": ids}, "status": 1}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("error decoding users")
	}

	return users, nil
}

func (r *userRepository) GetIdsExcluding(excludeIDs []primitive.ObjectID, limit int) ([]primitive.ObjectID, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	filter := bson.M{"status": 1}
	if len(excludeIDs) > 0 {
		filter["_id"] = bson.M{"$nin": excludeIDs}
	}

	pipeline := mongo.Pipeline{
		{{
			Key:   "$match",
			Value: filter,
		}},
		{{
			Key:   "$sample",
			Value: bson.M{"size": limit},
		}},
		{{
			Key:   "$project",
			Value: bson.M{"_id": 1},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error finding users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("error decoding users")
	}

	ids := make([]primitive.ObjectID, len(users))
	for i, user := range users {
		ids[i] = user.ID
	}

	return ids, nil
}
