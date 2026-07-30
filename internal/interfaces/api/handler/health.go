package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

var startTime = time.Now()

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	status := fiber.StatusOK
	checks := fiber.Map{
		"service": "finances-api",
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		status = fiber.StatusServiceUnavailable
		checks["database"] = "unhealthy"
	} else {
		checks["database"] = "healthy"
	}

	uptime := time.Since(startTime).String()
	checks["uptime"] = uptime

	return c.Status(status).JSON(fiber.Map{
		"success": status == fiber.StatusOK,
		"data":    checks,
	})
}
