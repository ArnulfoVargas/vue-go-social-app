package posts

import (
	"Server/internal/helpers"
	"Server/internal/store"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type postRepository struct {
	postCollection      *mongo.Collection
	postMediaCollection *mongo.Collection
}

func NewPostRepository(db *store.Database) *postRepository {
	return &postRepository{
		postCollection: db.Database.Collection("posts"),
	}
}

func (p *postRepository) CreatePost(post Post) error {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	_, err := p.postCollection.InsertOne(ctx, post)

	if err != nil {
		return errors.New("error creating post")
	}
	return nil
}

func (p *postRepository) GetPost(postId primitive.ObjectID) (Post, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	query := bson.M{"_id": postId}
	var post Post

	err := p.postCollection.FindOne(ctx, query).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Post{}, errors.New("post not found")
		}
		return Post{}, errors.New("error getting post")
	}

	return post, nil
}

func (p *postRepository) DeletePost(postId primitive.ObjectID) error {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	query := bson.M{"_id": postId}
	update := bson.M{"$set": bson.M{
		"status":    0,
		"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
	}}
	_, err := p.postCollection.UpdateOne(ctx, query, update)

	if err != nil {
		return errors.New("error deleting post")
	}
	return nil
}

func (p *postRepository) UpdatePost(postId primitive.ObjectID, update bson.M) (Post, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	updateQuery := bson.M{
		"$set":         update,
		"$currentDate": bson.M{"updatedAt": true},
	}

	_, err := p.postCollection.UpdateOne(ctx, bson.M{"_id": postId}, updateQuery)
	if err != nil {
		return Post{}, fmt.Errorf("updatePost: %w", err)
	}

	return p.GetPost(postId)
}

func (p *postRepository) GetPostsByUserId(userId primitive.ObjectID) ([]Post, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	cursor, err := p.postCollection.Find(ctx, bson.M{"userId": userId})
	if err != nil {
		return nil, errors.New("error getting posts")
	}
	defer cursor.Close(ctx)

	var posts = make([]Post, 0)
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, errors.New("error getting posts")
	}

	return posts, nil
}

func (p *postRepository) ExistsById(postId primitive.ObjectID) (bool, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	count, err := p.postCollection.CountDocuments(ctx, bson.M{"_id": postId})
	if err != nil {
		return false, errors.New("error checking post existence")
	}

	return count > 0, nil
}

func (p *postRepository) GetSuggestedPosts(userId primitive.ObjectID, limit int) ([]Post, error) {
	ctx, cancel := helpers.GenerateContext()
	defer cancel()

	pipeline := mongo.Pipeline{
		// Get all active accounts the user follows directly
		{{Key: "$match", Value: bson.M{
			"followerId": userId,
			"status":     1,
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":          nil,
			"followingIds": bson.M{"$push": "$followingId"},
		}}},

		// Find friends-of-friends (people your friends follow)
		{{Key: "$lookup", Value: bson.M{
			"from": "follows",
			"let":  bson.M{"followingIds": "$followingIds"},
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$in": bson.A{"$followerId", "$$followingIds"}},
						bson.M{"$eq": bson.A{"$status", 1}},
					},
				}}},
				bson.M{"$group": bson.M{
					"_id":              nil,
					"friendsFollowing": bson.M{"$push": "$followingId"},
				}},
			},
			"as": "fof",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$fof",
			"preserveNullAndEmptyArrays": true,
		}}},

		// Get authors of posts the user has previously liked
		{{Key: "$lookup", Value: bson.M{
			"from": "likes",
			"let":  bson.M{"userId": userId},
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$userId", "$$userId"}},
						bson.M{"$eq": bson.A{"$status", 1}},
					},
				}}},
				bson.M{"$lookup": bson.M{
					"from":         "posts",
					"localField":   "postId",
					"foreignField": "_id",
					"as":           "post",
				}},
				bson.M{"$unwind": "$post"},
				bson.M{"$group": bson.M{
					"_id":          nil,
					"likedAuthors": bson.M{"$push": "$post.userId"},
					"likedPostIds": bson.M{"$push": "$postId"},
				}},
			},
			"as": "likedData",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$likedData",
			"preserveNullAndEmptyArrays": true,
		}}},

		// Merge all relevant author sources into one array
		// Priority: direct follows > friends of friends > liked authors
		{{Key: "$project", Value: bson.M{
			"candidateAuthors": bson.M{
				"$concatArrays": bson.A{
					bson.M{"$ifNull": bson.A{"$followingIds", bson.A{}}},
					bson.M{"$ifNull": bson.A{"$fof.friendsFollowing", bson.A{}}},
					bson.M{"$ifNull": bson.A{"$likedData.likedAuthors", bson.A{}}},
				},
			},
			"likedPostIds": bson.M{"$ifNull": bson.A{"$likedData.likedPostIds", bson.A{}}},
		}}},

		// Fetch posts from all candidate authors
		{{Key: "$lookup", Value: bson.M{
			"from": "posts",
			"let": bson.M{
				"candidateAuthors": "$candidateAuthors",
				"likedPostIds":     "$likedPostIds",
			},
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"$expr": bson.M{
					"$and": bson.A{
						// Only posts from relevant authors
						bson.M{"$in": bson.A{"$userId", "$$candidateAuthors"}},
						// Exclude posts the user already liked
						bson.M{"$not": bson.A{bson.M{"$in": bson.A{"$_id", "$$likedPostIds"}}}},
						// Only active posts
						bson.M{"$eq": bson.A{"$status", 1}},
						// Exclude the user's own posts
						bson.M{"$ne": bson.A{"$userId", userId}},
					},
				}}},
				// Sort by most recent within the lookup
				bson.M{"$sort": bson.M{"createdAt": -1}},
			},
			"as": "posts",
		}}},

		{{Key: "$unwind", Value: "$posts"}},

		// Score each post based on signal strength
		// Direct follow = 3, friend of friend = 2, liked author = 1
		{{Key: "$addFields", Value: bson.M{
			"score": bson.M{
				"$add": bson.A{
					bson.M{"$cond": bson.A{
						bson.M{"$in": bson.A{"$posts.userId", bson.M{"$ifNull": bson.A{"$followingIds", bson.A{}}}}},
						3, 0,
					}},
					bson.M{"$cond": bson.A{
						bson.M{"$in": bson.A{"$posts.userId", bson.M{"$ifNull": bson.A{"$fof.friendsFollowing", bson.A{}}}}},
						2, 0,
					}},
					bson.M{"$cond": bson.A{
						bson.M{"$in": bson.A{"$posts.userId", bson.M{"$ifNull": bson.A{"$likedData.likedAuthors", bson.A{}}}}},
						1, 0,
					}},
				},
			},
		}}},

		// Sort by score then recency, return top N
		{{Key: "$sort", Value: bson.D{
			{Key: "score", Value: -1},
			{Key: "posts.createdAt", Value: -1},
		}}},
		{{Key: "$limit", Value: limit}},

		// Reshape output to match model.Post
		{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$posts"}}},
	}

	cursor, err := p.postCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("getRecommendedPosts aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	posts := make([]Post, 0)
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("getRecommendedPosts decode: %w", err)
	}

	return posts, nil
}
