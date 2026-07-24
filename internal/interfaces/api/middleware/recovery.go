package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/agnathor/finances-go/pkg/logger"
	"github.com/agnathor/finances-go/pkg/response"
)

func Recovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				logger.Get().Error("panic recovered",
					zap.Any("panic", r),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
				)
				_ = response.Error(c, fiber.StatusInternalServerError, "internal server error")
			}
		}()
		return c.Next()
	}
}
