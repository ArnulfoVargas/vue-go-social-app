package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Like struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	PostID    primitive.ObjectID `bson:"postId" json:"postId" validate:"required"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId" validate:"required"`
	Status    uint8              `bson:"status" json:"status" validate:"required"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
}
