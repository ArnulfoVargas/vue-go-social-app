package users

import (
	"Server/internal/features/media"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Bio       string             `json:"bio" bson:"bio"`
	Avatar    *media.Media       `json:"avatar,omitempty" bson:"avatar,omitempty"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
	Status    uint8              `json:"status" bson:"status"`
}
