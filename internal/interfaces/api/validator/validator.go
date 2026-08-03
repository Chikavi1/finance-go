package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("iso_currency", func(fl validator.FieldLevel) bool {
		switch strings.ToUpper(fl.Field().String()) {
		case "USD", "EUR", "BRL", "GBP", "JPY", "MXN":
			return true
		}
		return false
	})
}

func ValidateRequest(c *fiber.Ctx, req interface{}) error {
	if err := c.BodyParser(req); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}

	if err := validate.Struct(req); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			errors[field] = fmt.Sprintf(
				"failed %s validation (value: %v)",
				err.Tag(),
				err.Value(),
			)
		}
		return &ValidationError{Errors: errors}
	}

	return nil
}

type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}
