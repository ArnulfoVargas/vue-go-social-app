package repository

import (
	"Server/internal/helpers"
	"Server/internal/model"
	"Server/internal/store"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type followRepository struct {
	collection *mongo.Collection
}

func NewFollowRepository(db *store.Database) *followRepository {
	return &followRepository{
		collection: db.Database.Collection("follows"),
	}
}

func (r *followRepository) FollowUser(userID, targetUserID string) error {
	exists, err := r.existsFollowUnscoped(userID, targetUserID)

	if err != nil {
		return err
	}

	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	userId, err := helpers.ToObjectID(userID)
	if err != nil {
		return err
	}
	targetId, err := helpers.ToObjectID(targetUserID)
	if err != nil {
		return err
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	if exists {
		data := bson.M{
			"status":    1,
			"updatedAt": now,
		}
		_, err := r.collection.UpdateOne(ctx, bson.M{"followerId": userId, "followingId": targetId}, bson.M{"$set": data})
		if err != nil {
			return fmt.Errorf("error creating follow")
		}
		return nil
	}

	follow := model.Follow{
		ID:          primitive.NewObjectID(),
		FollowerID:  userId,
		FollowingID: targetId,
		Status:      1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err = r.collection.InsertOne(ctx, follow)
	if err != nil {
		return fmt.Errorf("error creating follow")
	}
	return nil
}

func (r *followRepository) UnfollowUser(userID, targetUserID string) error {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	exists, err := r.existsFollowUnscoped(userID, targetUserID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("cannot unfollow")
	}

	userId, err := helpers.ToObjectID(userID)
	if err != nil {
		return err
	}
	targetId, err := helpers.ToObjectID(targetUserID)
	if err != nil {
		return err
	}

	now := time.Now()
	data := bson.M{
		"status":    0,
		"updatedAt": primitive.NewDateTimeFromTime(now),
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"followerId": userId, "followingId": targetId}, bson.M{"$set": data})
	if err != nil {
		return fmt.Errorf("error unfollowing user")
	}
	return nil
}

func (r *followRepository) UserIsFollowing(userID, targetUserID string) (bool, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	idUser, err := helpers.ToObjectID(userID)
	if err != nil {
		return false, err
	}
	idTarget, err := helpers.ToObjectID(targetUserID)
	if err != nil {
		return false, err
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"followerId": idUser, "followingId": idTarget, "status": 1})
	if err != nil {
		return false, fmt.Errorf("error checking follow status")
	}

	return count > 0, nil
}

func (r *followRepository) existsFollowUnscoped(userId, targetId string) (bool, error) {
	col := r.collection
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	idUser, err := helpers.ToObjectID(userId)
	if err != nil {
		return false, err
	}
	idTarget, err := helpers.ToObjectID(targetId)
	if err != nil {
		return false, err
	}

	count, err := col.CountDocuments(ctx, bson.M{"followerId": idUser, "followingId": idTarget})
	if err != nil {
		return false, fmt.Errorf("error searching follow")
	}

	return count > 0, nil
}

func (r *followRepository) GetFollowingCount(userID string) (int64, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	idUser, err := helpers.ToObjectID(userID)
	if err != nil {
		return 0, err
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"followerId": idUser, "status": 1})
	if err != nil {
		return 0, fmt.Errorf("error getting following count")
	}

	return count, nil
}

func (r *followRepository) GetFollowerCount(userID string) (int64, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	idUser, err := helpers.ToObjectID(userID)
	if err != nil {
		return 0, err
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"followingId": idUser, "status": 1})
	if err != nil {
		return 0, fmt.Errorf("error getting follower count")
	}

	return count, nil
}

func (r *followRepository) GetFollowingIds(userID string) ([]primitive.ObjectID, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	idUser, err := helpers.ToObjectID(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx,
		bson.M{"followerId": idUser, "status": 1},
		options.Find().SetProjection(bson.M{"followingId": 1, "_id": 0}))

	if err != nil {
		return nil, fmt.Errorf("error getting following ids")
	}
	defer cursor.Close(ctx)

	var results []struct {
		FollowingId primitive.ObjectID `bson:"followingId"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("error decoding following ids")
	}

	followingIds := make([]primitive.ObjectID, len(results))
	for i, result := range results {
		followingIds[i] = result.FollowingId
	}

	return followingIds, nil
}

func (r *followRepository) GetRelatedFollowSuggestions(userId primitive.ObjectID, followingIds []primitive.ObjectID, limit int) ([]primitive.ObjectID, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	pipeline := mongo.Pipeline{
		// Get all follow relationships where the user is either a follower or following
		{{Key: "$match", Value: bson.M{
			"$or": bson.A{
				bson.M{"followerId": bson.M{"$in": followingIds}, "status": 1},
				bson.M{"followingId": bson.M{"$in": followingIds}, "status": 1},
			},
		}}},
		// Run both directions in parallel
		{{Key: "$facet", Value: bson.M{
			// People who your friends follow
			"following": bson.A{
				bson.M{"$match": bson.M{"followerId": bson.M{"$in": followingIds}}},
				bson.M{"$project": bson.M{"userId": "$followingId", "_id": 0}},
			},
			// People who follow your friends
			"followers": bson.A{
				bson.M{"$match": bson.M{"followingId": bson.M{"$in": followingIds}}},
				bson.M{"$project": bson.M{"userId": "$followerId", "_id": 0}},
			},
		}}},
		// Merge both arrays into one
		{{Key: "$project", Value: bson.M{
			"users": bson.M{"$concatArrays": bson.A{"$following", "$followers"}},
		}}},
		{{Key: "$unwind", Value: "$users"}},
		// Exclude yourself and people you already follow
		{{Key: "$match", Value: bson.M{
			"users.userId": bson.M{"$nin": append(followingIds, userId)},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$users.userId",
			"score": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"score": -1}}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error getting related follow suggestions")
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("error decoding related follow suggestions")
	}

	suggestedIds := make([]primitive.ObjectID, len(results))
	for i, result := range results {
		suggestedIds[i] = result.ID
	}

	return suggestedIds, nil
}
