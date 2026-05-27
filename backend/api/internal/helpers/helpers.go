package helpers

import (
	"Server/internal/constants"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToObjectID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid id: %s", id)
	}
	return oid, nil
}

func GenerateContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), constants.DEFAULT_CONTEXT_TIMEOUT)
}

func GetUserIdFromLocals(c fiber.Ctx) (string, bool) {
	id, ok := c.Locals(constants.USER_ID_CLAIM).(string)
	return id, ok
}
