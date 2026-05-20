package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	PostID    primitive.ObjectID `bson:"postId" json:"postId"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
}
