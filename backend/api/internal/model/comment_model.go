package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	PostID    primitive.ObjectID `bson:"postId" json:"postId" validate:"required"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId" validate:"required"`
	Content   string             `bson:"content" json:"content" validate:"min=0,max=500"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
}
