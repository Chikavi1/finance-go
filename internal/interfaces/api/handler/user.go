package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/internal/application/user"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/internal/interfaces/api/validator"
	"github.com/agnathor/finances-go/pkg/response"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	user, err := h.userService.GetProfile(c.Context(), userID)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "user not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to get profile")
	}

	return response.JSON(c, fiber.StatusOK, mapUserToResponse(user))
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req dto.UpdateProfileRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	user, err := h.userService.UpdateProfile(c.Context(), userID, req.Name, req.AvatarURL)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to update profile")
	}

	return response.JSON(c, fiber.StatusOK, mapUserToResponse(user))
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req dto.ChangePasswordRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.userService.ChangePassword(c.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		if err == domain.ErrInvalidCredentials {
			return response.Error(c, fiber.StatusUnauthorized, "current password is incorrect")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to change password")
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"message": "password changed successfully",
	})
}
