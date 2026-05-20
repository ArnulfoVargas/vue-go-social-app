package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Content   string             `bson:"content,omitempty" json:"content,omitempty"`
	Status    uint8              `bson:"status,default=1" json:"status,omitempty"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
}

type PostMedia struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	PostID primitive.ObjectID `bson:"postId" json:"postId"`
	Media  []Media            `bson:"media" json:"media"`
}
