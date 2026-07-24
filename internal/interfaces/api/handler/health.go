package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"status":  "ok",
		"service": "finances-api",
	})
}
