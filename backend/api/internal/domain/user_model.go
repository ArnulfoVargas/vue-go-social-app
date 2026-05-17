package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string               `json:"name" bson:"name" validate:"required,min=3,max=100"`
	Email     string               `json:"email" bson:"email" validate:"required,email"`
	Password  string               `json:"password" bson:"password" validate:"required,min=8,max=32"`
	ImageUrl  string               `json:"imageUrl,omitempty" bson:"imageUrl,omitempty" validate:"url"`
	Bio       string               `json:"bio" bson:"bio" validate:"max=255"`
	Followers []primitive.ObjectID `json:"followers" bson:"followers"`
	Following []primitive.ObjectID `json:"following" bson:"following"`
	CreatedAt primitive.DateTime   `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime   `json:"updatedAt" bson:"updatedAt"`
	Status    uint8                `json:"status" bson:"status"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
