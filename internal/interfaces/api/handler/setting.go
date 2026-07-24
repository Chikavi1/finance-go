package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/internal/application/setting"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/internal/interfaces/api/validator"
	"github.com/agnathor/finances-go/pkg/response"
)

type SettingHandler struct {
	settingService setting.Service
}

func NewSettingHandler(settingService setting.Service) *SettingHandler {
	return &SettingHandler{settingService: settingService}
}

func (h *SettingHandler) GetAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	settings, err := h.settingService.GetAll(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get settings")
	}

	resp := make(map[string]string, len(settings))
	for _, s := range settings {
		resp[s.Key] = s.Value
	}

	return response.JSON(c, fiber.StatusOK, resp)
}

func (h *SettingHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req dto.UpdateSettingsRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	settings, err := h.settingService.Update(c.Context(), userID, req.Settings)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to update settings")
	}

	resp := make(map[string]string, len(settings))
	for _, s := range settings {
		resp[s.Key] = s.Value
	}

	return response.JSON(c, fiber.StatusOK, resp)
}
