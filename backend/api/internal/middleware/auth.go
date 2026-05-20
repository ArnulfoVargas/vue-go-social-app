package middleware

import (
	"Server/internal/constants"
	"Server/internal/service"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func Protected(parseJWT func(string) (*service.Claims, error)) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := parseJWT(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		c.Locals(constants.USER_ID_CLAIM, claims.UserID)
		c.Locals(constants.MAIL_CLAIM, claims.Email)
		return c.Next()
	}
}
