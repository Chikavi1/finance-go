package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/internal/application/budget"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/internal/interfaces/api/validator"
	"github.com/agnathor/finances-go/pkg/response"
)

type BudgetHandler struct {
	budgetService budget.Service
}

func NewBudgetHandler(budgetService budget.Service) *BudgetHandler {
	return &BudgetHandler{budgetService: budgetService}
}

func (h *BudgetHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req dto.CreateBudgetRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	budget, err := h.budgetService.Create(c.Context(), userID, req.CategoryID, req.Amount, req.Month, req.Year)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create budget")
	}

	return response.Created(c, mapBudgetToResponse(budget))
}

func (h *BudgetHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	budgetID := c.Params("id")

	budget, err := h.budgetService.GetByID(c.Context(), userID, budgetID)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "budget not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to get budget")
	}

	return response.JSON(c, fiber.StatusOK, mapBudgetToResponse(budget))
}

func (h *BudgetHandler) GetAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	month := c.QueryInt("month", 0)
	year := c.QueryInt("year", 0)

	var budgets []*domain.Budget
	var err error

	if month > 0 && year > 0 {
		budgets, err = h.budgetService.GetByMonthYear(c.Context(), userID, int32(month), int32(year))
	} else {
		budgets, err = h.budgetService.GetAll(c.Context(), userID)
	}

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get budgets")
	}

	resp := make([]dto.BudgetResponse, len(budgets))
	for i, b := range budgets {
		resp[i] = mapBudgetToResponse(b)
	}

	return response.JSON(c, fiber.StatusOK, resp)
}

func (h *BudgetHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	budgetID := c.Params("id")

	var req dto.UpdateBudgetRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	budget, err := h.budgetService.Update(c.Context(), userID, budgetID, req.Amount, req.Spent)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "budget not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to update budget")
	}

	return response.JSON(c, fiber.StatusOK, mapBudgetToResponse(budget))
}

func (h *BudgetHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	budgetID := c.Params("id")

	if err := h.budgetService.Delete(c.Context(), userID, budgetID); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "budget not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete budget")
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"message": "budget deleted successfully",
	})
}

func mapBudgetToResponse(b *domain.Budget) dto.BudgetResponse {
	return dto.BudgetResponse{
		ID:         b.ID,
		UserID:     b.UserID,
		CategoryID: b.CategoryID,
		Amount:     b.Amount,
		Spent:      b.Spent,
		Month:      b.Month,
		Year:       b.Year,
		CreatedAt:  b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  b.UpdatedAt.Format(time.RFC3339),
	}
}
