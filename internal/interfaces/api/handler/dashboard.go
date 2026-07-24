package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/internal/application/dashboard"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/pkg/response"
)

type DashboardHandler struct {
	dashboardService dashboard.Service
}

func NewDashboardHandler(dashboardService dashboard.Service) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	data, err := h.dashboardService.GetDashboard(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get dashboard data")
	}

	txns := make([]dto.TransactionResponse, len(data.RecentTxns))
	for i, tx := range data.RecentTxns {
		txns[i] = mapTransactionToResponse(tx)
	}

	budgets := make([]dto.BudgetResponse, len(data.Budgets))
	for i, b := range data.Budgets {
		budgets[i] = mapBudgetToResponse(b)
	}

	resp := dto.DashboardResponse{
		TotalBalance:     data.TotalBalance,
		MonthlyIncome:    data.MonthlyIncome,
		MonthlyExpenses:  data.MonthlyExpenses,
		RecentTxns:       txns,
		Budgets:          budgets,
		ActiveDebtsTotal: data.ActiveDebtsTotal,
		ActiveDebtCount:  data.ActiveDebtCount,
	}

	return response.JSON(c, fiber.StatusOK, resp)
}
