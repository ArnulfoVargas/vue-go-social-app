package likes

import (
	"Server/internal/helpers"
	"Server/internal/store"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type likeRepository struct {
	collection *mongo.Collection
}

func NewlikeRepository(db *store.Database) *likeRepository {
	return &likeRepository{collection: db.Database.Collection("likes")}
}

func (r *likeRepository) DeleteLike(like Like) error {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": like.ID}, bson.M{"$set": bson.M{
		"status": 0,
	}})
	if err != nil {
		return errors.New("cannot delete")
	}
	if result.ModifiedCount == 0 {
		return errors.New("no document deleted")
	}
	return nil
}

func (r *likeRepository) AddLike(like Like) error {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	if hasLike, err := r.hasLikeUnscoped(like.PostID, like.UserID); err != nil {
		return err
	} else if hasLike {
		return r.LikePost(like)
	}
	_, err := r.collection.InsertOne(ctx, like)
	if err != nil {
		return errors.New("cannot add like")
	}
	return nil
}

func (r *likeRepository) LikePost(like Like) error {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	result, err := r.collection.UpdateOne(ctx, bson.M{"userId": like.UserID, "postId": like.PostID}, bson.M{"$set": bson.M{
		"status": 1,
	}})
	if err != nil {
		return errors.New("cannot like post")
	}
	if result.ModifiedCount == 0 {
		return errors.New("no document updated")
	}

	return nil
}

func (r *likeRepository) DeleteLikesFromPost(postId primitive.ObjectID) error {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	_, err := r.collection.UpdateMany(ctx, bson.M{"postId": postId}, bson.M{"$set": bson.M{
		"status": 0,
	}})
	if err != nil {
		return errors.New("cannot delete likes from post")
	}
	return nil
}

func (r *likeRepository) HasLike(postId, userId primitive.ObjectID) (bool, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"postId": postId, "userId": userId, "status": 1})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *likeRepository) hasLikeUnscoped(postId, userId primitive.ObjectID) (bool, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"postId": postId, "userId": userId})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
