package posts

import (
	"Server/internal/features/media"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Content   string             `bson:"content,omitempty" json:"content,omitempty"`
	Status    uint8              `bson:"status,default=1" json:"status,omitempty"`
	Media     []media.Media      `bson:"media,omitempty" json:"media,omitempty"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
}
