package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/internal/application/account"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/internal/interfaces/api/validator"
	"github.com/agnathor/finances-go/pkg/response"
)

type AccountHandler struct {
	accountService account.Service
}

func NewAccountHandler(accountService account.Service) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

func (h *AccountHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req dto.CreateAccountRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	account, err := h.accountService.Create(c.Context(), userID, req.Name, req.Type, req.Currency, req.Balance, req.Color, req.Icon)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create account")
	}

	return response.Created(c, mapAccountToResponse(account))
}

func (h *AccountHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	accountID := c.Params("id")

	account, err := h.accountService.GetByID(c.Context(), userID, accountID)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "account not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to get account")
	}

	return response.JSON(c, fiber.StatusOK, mapAccountToResponse(account))
}

func (h *AccountHandler) GetAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	accounts, err := h.accountService.GetAll(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get accounts")
	}

	resp := make([]dto.AccountResponse, len(accounts))
	for i, a := range accounts {
		resp[i] = mapAccountToResponse(a)
	}

	return response.JSON(c, fiber.StatusOK, resp)
}

func (h *AccountHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	accountID := c.Params("id")

	var req dto.UpdateAccountRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	account, err := h.accountService.Update(c.Context(), userID, accountID, req.Name, req.Type, req.Currency, req.Color, req.Icon, req.Archived)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "account not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to update account")
	}

	return response.JSON(c, fiber.StatusOK, mapAccountToResponse(account))
}

func (h *AccountHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	accountID := c.Params("id")

	if err := h.accountService.Delete(c.Context(), userID, accountID); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "account not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete account")
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"message": "account deleted successfully",
	})
}

func mapAccountToResponse(account *domain.Account) dto.AccountResponse {
	return dto.AccountResponse{
		ID:        account.ID,
		UserID:    account.UserID,
		Name:      account.Name,
		Type:      account.Type,
		Currency:  account.Currency,
		Balance:   account.Balance,
		Color:     account.Color,
		Icon:      account.Icon,
		Archived:  account.Archived,
		CreatedAt: account.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: account.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
