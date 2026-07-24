package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/pkg/jwt"
	"github.com/agnathor/finances-go/pkg/response"
)

func AuthRequired(jwtManager *jwt.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, fiber.StatusUnauthorized, "invalid authorization header format")
		}

		claims, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "invalid or expired token")
		}

		c.Locals("user_id", claims.UserID)
		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) string {
	if userID, ok := c.Locals("user_id").(string); ok {
		return userID
	}
	return ""
}
