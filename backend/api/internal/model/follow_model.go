package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Follow struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FollowerID  primitive.ObjectID `bson:"followerId" json:"followerId"`
	FollowingID primitive.ObjectID `bson:"followingId" json:"followingId"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
	Status      uint8              `json:"status" bson:"status"`
}
