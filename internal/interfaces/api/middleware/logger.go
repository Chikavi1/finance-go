package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/agnathor/finances-go/pkg/logger"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		status := c.Response().StatusCode()
		reqLogger := logger.Get()

		fields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		}

		if userID := c.Locals("user_id"); userID != nil {
			fields = append(fields, zap.String("user_id", userID.(string)))
		}

		switch {
		case status >= 500:
			reqLogger.Error("request failed", fields...)
		case status >= 400:
			reqLogger.Warn("request warning", fields...)
		default:
			reqLogger.Info("request completed", fields...)
		}

		return err
	}
}
