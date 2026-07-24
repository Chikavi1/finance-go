package response

import (
	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
	Total   int `json:"total,omitempty"`
}

func JSON(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(APIResponse{
		Success: true,
		Data:    data,
	})
}

func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success: true,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(APIResponse{
		Success: false,
		Error:   message,
	})
}

func ValidationError(c *fiber.Ctx, errors map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"success": false,
		"error":   "validation failed",
		"fields":  errors,
	})
}

func Paginated(c *fiber.Ctx, data interface{}, page, perPage, total int) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:    page,
			PerPage: perPage,
			Total:   total,
		},
	})
}
