package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/internal/application/auth"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/validator"
	"github.com/agnathor/finances-go/pkg/response"
)

type AuthHandler struct {
	authService auth.Service
}

func NewAuthHandler(authService auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	user, tokens, err := h.authService.Register(c.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		if err == domain.ErrConflict {
			return response.Error(c, fiber.StatusConflict, "email already registered")
		}
		return response.Error(c, fiber.StatusInternalServerError, "registration failed")
	}

	return response.Created(c, dto.AuthResponse{
		User:         mapUserToResponse(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	user, tokens, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			return response.Error(c, fiber.StatusUnauthorized, "invalid email or password")
		}
		return response.Error(c, fiber.StatusInternalServerError, "login failed")
	}

	return response.JSON(c, fiber.StatusOK, dto.AuthResponse{
		User:         mapUserToResponse(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	tokens, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		if err == domain.ErrTokenInvalid || err == domain.ErrTokenExpired {
			return response.Error(c, fiber.StatusUnauthorized, "invalid or expired refresh token")
		}
		return response.Error(c, fiber.StatusInternalServerError, "token refresh failed")
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_at":    tokens.ExpiresAt,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req dto.LogoutRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.authService.Logout(c.Context(), userID, req.RefreshToken); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "logout failed")
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"message": "logged out successfully",
	})
}

func mapUserToResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
