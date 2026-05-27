package media

import "go.mongodb.org/mongo-driver/bson/primitive"

type Media struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	URL      string             `bson:"url" json:"url"`
	PublicID string             `bson:"publicId" json:"publicId"`
}
