package repository

import (
	"Server/internal/helpers"
	"Server/internal/model"
	"Server/internal/store"
	"fmt"

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

func (r *followRepository) FollowUser(follow model.Follow) error {
	exists, err := r.existsFollowUnscoped(follow.FollowerID, follow.FollowingID)

	if err != nil {
		return err
	}

	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	if exists {
		query := bson.M{"$set": bson.M{"status": 1},
			"$currentDate": bson.M{"updatedAt": true}}

		_, err := r.collection.UpdateOne(ctx, bson.M{"followerId": follow.FollowerID, "followingId": follow.FollowingID}, query)
		if err != nil {
			return fmt.Errorf("error creating follow")
		}
		return nil
	}

	_, err = r.collection.InsertOne(ctx, follow)
	if err != nil {
		return fmt.Errorf("error creating follow")
	}
	return nil
}

func (r *followRepository) UnfollowUser(userID, targetUserID primitive.ObjectID) error {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	exists, err := r.existsFollowUnscoped(userID, targetUserID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("cannot unfollow")
	}

	query := bson.M{"$set": bson.M{"status": 0},
		"$currentDate": bson.M{"updatedAt": true},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"followerId": userID, "followingId": targetUserID}, query)
	if err != nil {
		return fmt.Errorf("error unfollowing user")
	}
	return nil
}

func (r *followRepository) UserIsFollowing(userID, targetUserID primitive.ObjectID) (bool, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"followerId": userID, "followingId": targetUserID, "status": 1})
	if err != nil {
		return false, fmt.Errorf("error checking follow status")
	}

	return count > 0, nil
}

func (r *followRepository) existsFollowUnscoped(userId, targetId primitive.ObjectID) (bool, error) {
	col := r.collection
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	count, err := col.CountDocuments(ctx, bson.M{"followerId": userId, "followingId": targetId})
	if err != nil {
		return false, fmt.Errorf("error searching follow")
	}

	return count > 0, nil
}

func (r *followRepository) GetFollowingCount(userID primitive.ObjectID) (int64, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"followerId": userID, "status": 1})
	if err != nil {
		return 0, fmt.Errorf("error getting following count")
	}

	return count, nil
}

func (r *followRepository) GetFollowerCount(userID primitive.ObjectID) (int64, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"followingId": userID, "status": 1})
	if err != nil {
		return 0, fmt.Errorf("error getting follower count")
	}

	return count, nil
}

func (r *followRepository) GetFollowingIds(userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	cursor, err := r.collection.Find(ctx,
		bson.M{"followerId": userID, "status": 1},
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

	excludedIds := make([]primitive.ObjectID, len(followingIds)+1)
	copy(excludedIds, followingIds)
	excludedIds[len(followingIds)] = userId

	pipeline := mongo.Pipeline{
		// Match all follow relationships where any of the user's friends
		// are involved (either as follower or following)
		{{Key: "$match", Value: bson.M{
			"$or": bson.A{
				bson.M{"followerId": bson.M{"$in": followingIds}, "status": 1},
				bson.M{"followingId": bson.M{"$in": followingIds}, "status": 1},
			},
		}}},
		// Split into two branches in parallel:
		// - "following": people that your friends follow
		// - "followers": people that follow your friends
		{{Key: "$facet", Value: bson.M{
			"following": bson.A{
				bson.M{"$match": bson.M{"followerId": bson.M{"$in": followingIds}, "status": 1}},
				bson.M{"$project": bson.M{"users": bson.M{"userId": "$followingId"}, "_id": 0}},
			},
			"followers": bson.A{
				bson.M{"$match": bson.M{"followingId": bson.M{"$in": followingIds}, "status": 1}},
				bson.M{"$project": bson.M{"users": bson.M{"userId": "$followerId"}, "_id": 0}},
			},
		}}},
		// Merge both arrays into a single "users" array
		{{Key: "$project", Value: bson.M{
			"users": bson.M{"$concatArrays": bson.A{"$following", "$followers"}},
		}}},
		// Flatten the users array so each userId becomes its own document
		{{Key: "$unwind", Value: "$users"}},
		// Remove the current user and people they already follow
		// to avoid suggesting someone they are already connected to
		{{Key: "$match", Value: bson.M{
			"users.userId": bson.M{"$nin": excludedIds},
		}}},
		// Inject an additional signal: owners of posts the user has liked.
		// Each liked post owner gets +1 to their score, same as a follow connection.
		// This rewards accounts the user actively engages with.
		{{Key: "$unionWith", Value: bson.M{
			"coll": "likes",
			"pipeline": bson.A{
				// Only consider likes made by the current user
				bson.M{"$match": bson.M{"userId": userId}},
				// Look up the post to get the owner's userId
				bson.M{"$lookup": bson.M{
					"from":         "posts",
					"localField":   "postId",
					"foreignField": "_id",
					"as":           "post",
				}},
				bson.M{"$unwind": "$post"},
				// Normalize to the same {users: {userId}} shape as the follow pipeline
				bson.M{"$project": bson.M{
					"users": bson.M{"userId": "$post.userId"},
					"_id":   0,
				}},
				bson.M{"$unwind": "$users"},
				// Apply the same exclusion filter as the follow pipeline
				bson.M{"$match": bson.M{
					"users.userId": bson.M{"$nin": excludedIds},
				}},
			},
		}}},
		// Count how many times each candidate userId appeared across
		// both follow connections and liked posts — this is their relevance score
		{{Key: "$group", Value: bson.M{
			"_id":   "$users.userId",
			"score": bson.M{"$sum": 1},
		}}},
		// Hydrate each userId with their full user document
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "user",
		}}},
		{{Key: "$unwind", Value: "$user"}},
		// Filter out inactive accounts
		{{Key: "$match", Value: bson.M{"user.status": 1}}},
		// Return the most relevant suggestions first
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
