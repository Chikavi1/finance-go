package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/internal/application/transaction"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/internal/interfaces/api/validator"
	"github.com/agnathor/finances-go/pkg/response"
)

type TransactionHandler struct {
	transactionService transaction.Service
}

func NewTransactionHandler(transactionService transaction.Service) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req dto.CreateTransactionRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
	}

	tx := &domain.Transaction{
		AccountID:   req.AccountID,
		ToAccountID: req.ToAccountID,
		CategoryID:  req.CategoryID,
		Type:        domain.TransactionType(req.Type),
		Amount:      req.Amount,
		Description: req.Description,
		Notes:       req.Notes,
		Date:        date,
		Tags:        req.Tags,
	}

	if err := h.transactionService.Create(c.Context(), userID, tx); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create transaction")
	}

	return response.Created(c, mapTransactionToResponse(tx))
}

func (h *TransactionHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id := c.Params("id")

	tx, err := h.transactionService.GetByID(c.Context(), userID, id)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "transaction not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to get transaction")
	}

	return response.JSON(c, fiber.StatusOK, mapTransactionToResponse(tx))
}

func (h *TransactionHandler) GetAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var filter domain.TransactionFilter

	if t := c.Query("type"); t != "" {
		tType := domain.TransactionType(t)
		filter.Type = &tType
	}

	if a := c.Query("account_id"); a != "" {
		filter.AccountID = &a
	}

	if s := c.Query("start_date"); s != "" {
		start, err := time.Parse("2006-01-02", s)
		if err == nil {
			filter.StartDate = &start
		}
	}

	if e := c.Query("end_date"); e != "" {
		end, err := time.Parse("2006-01-02", e)
		if err == nil {
			filter.EndDate = &end
		}
	}

	transactions, err := h.transactionService.GetAll(c.Context(), userID, filter)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get transactions")
	}

	responses := make([]dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = mapTransactionToResponse(tx)
	}

	return response.JSON(c, fiber.StatusOK, responses)
}

func (h *TransactionHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id := c.Params("id")

	var req dto.UpdateTransactionRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
	}

	tx := &domain.Transaction{
		ID:          id,
		AccountID:   req.AccountID,
		ToAccountID: req.ToAccountID,
		CategoryID:  req.CategoryID,
		Type:        domain.TransactionType(req.Type),
		Amount:      req.Amount,
		Description: req.Description,
		Notes:       req.Notes,
		Date:        date,
		Tags:        req.Tags,
	}

	if err := h.transactionService.Update(c.Context(), userID, tx); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "transaction not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to update transaction")
	}

	return response.JSON(c, fiber.StatusOK, mapTransactionToResponse(tx))
}

func (h *TransactionHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	id := c.Params("id")

	if err := h.transactionService.Delete(c.Context(), userID, id); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "transaction not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete transaction")
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"message": "transaction deleted successfully",
	})
}

func mapTransactionToResponse(tx *domain.Transaction) dto.TransactionResponse {
	return dto.TransactionResponse{
		ID:          tx.ID,
		AccountID:   tx.AccountID,
		ToAccountID: tx.ToAccountID,
		CategoryID:  tx.CategoryID,
		Type:        string(tx.Type),
		Amount:      tx.Amount,
		Description: tx.Description,
		Notes:       tx.Notes,
		Date:        tx.Date.Format("2006-01-02"),
		Tags:        tx.Tags,
		CreatedAt:   tx.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   tx.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
