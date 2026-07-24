package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	debtService "github.com/agnathor/finances-go/internal/application/debt"
	debtPaymentService "github.com/agnathor/finances-go/internal/application/debtpayment"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/internal/interfaces/api/validator"
	"github.com/agnathor/finances-go/pkg/response"
)

type DebtHandler struct {
	debtService        debtService.Service
	debtPaymentService debtPaymentService.Service
}

func NewDebtHandler(debtService debtService.Service, debtPaymentService debtPaymentService.Service) *DebtHandler {
	return &DebtHandler{debtService: debtService, debtPaymentService: debtPaymentService}
}

func (h *DebtHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req dto.CreateDebtRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid due_date format, use YYYY-MM-DD")
		}
		dueDate = &t
	}

	debt, err := h.debtService.Create(c.Context(), userID, req.Name, req.TotalAmount, req.RemainingAmount, req.InterestRate, dueDate, domain.DebtStatus(req.Status), req.Notes)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create debt")
	}

	return response.Created(c, mapDebtToResponse(debt))
}

func (h *DebtHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	debtID := c.Params("id")

	debt, err := h.debtService.GetByID(c.Context(), userID, debtID)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "debt not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to get debt")
	}

	return response.JSON(c, fiber.StatusOK, mapDebtToResponse(debt))
}

func (h *DebtHandler) GetAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	debts, err := h.debtService.GetAll(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get debts")
	}

	resp := make([]dto.DebtResponse, len(debts))
	for i, d := range debts {
		resp[i] = mapDebtToResponse(d)
	}

	return response.JSON(c, fiber.StatusOK, resp)
}

func (h *DebtHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	debtID := c.Params("id")

	var req dto.UpdateDebtRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid due_date format, use YYYY-MM-DD")
		}
		dueDate = &t
	}

	debt, err := h.debtService.Update(c.Context(), userID, debtID, req.Name, req.TotalAmount, req.RemainingAmount, req.InterestRate, dueDate, domain.DebtStatus(req.Status), req.Notes)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "debt not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to update debt")
	}

	return response.JSON(c, fiber.StatusOK, mapDebtToResponse(debt))
}

func (h *DebtHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	debtID := c.Params("id")

	if err := h.debtService.Delete(c.Context(), userID, debtID); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "debt not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete debt")
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"message": "debt deleted successfully",
	})
}

func (h *DebtHandler) CreatePayment(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	debtID := c.Params("id")

	_, err := h.debtService.GetByID(c.Context(), userID, debtID)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "debt not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to get debt")
	}

	var req dto.CreateDebtPaymentRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid payment_date format, use YYYY-MM-DD")
	}

	payment, err := h.debtPaymentService.Create(c.Context(), debtID, req.Amount, paymentDate, req.Notes)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create debt payment")
	}

	return response.Created(c, mapDebtPaymentToResponse(payment))
}

func (h *DebtHandler) GetPayments(c *fiber.Ctx) error {
	debtID := c.Params("id")

	payments, err := h.debtPaymentService.GetByDebtID(c.Context(), debtID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get debt payments")
	}

	resp := make([]dto.DebtPaymentResponse, len(payments))
	for i, p := range payments {
		resp[i] = mapDebtPaymentToResponse(p)
	}

	return response.JSON(c, fiber.StatusOK, resp)
}

func (h *DebtHandler) DeletePayment(c *fiber.Ctx) error {
	paymentID := c.Params("paymentId")

	if err := h.debtPaymentService.Delete(c.Context(), paymentID); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "debt payment not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete debt payment")
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"message": "debt payment deleted successfully",
	})
}

func mapDebtToResponse(d *domain.Debt) dto.DebtResponse {
	var dueDate *string
	if d.DueDate != nil {
		t := d.DueDate.Format("2006-01-02")
		dueDate = &t
	}

	return dto.DebtResponse{
		ID:              d.ID,
		UserID:          d.UserID,
		Name:            d.Name,
		TotalAmount:     d.TotalAmount,
		RemainingAmount: d.RemainingAmount,
		InterestRate:    d.InterestRate,
		DueDate:         dueDate,
		Status:          string(d.Status),
		Notes:           d.Notes,
		CreatedAt:       d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       d.UpdatedAt.Format(time.RFC3339),
	}
}

func mapDebtPaymentToResponse(p *domain.DebtPayment) dto.DebtPaymentResponse {
	return dto.DebtPaymentResponse{
		ID:          p.ID,
		DebtID:      p.DebtID,
		Amount:      p.Amount,
		PaymentDate: p.PaymentDate.Format("2006-01-02"),
		Notes:       p.Notes,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
	}
}
