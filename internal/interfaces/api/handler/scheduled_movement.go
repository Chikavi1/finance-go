package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	scheduledMovementService "github.com/agnathor/finances-go/internal/application/scheduledmovement"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/internal/interfaces/api/validator"
	"github.com/agnathor/finances-go/pkg/response"
)

type ScheduledMovementHandler struct {
	scheduledMovementService scheduledMovementService.Service
}

func NewScheduledMovementHandler(scheduledMovementService scheduledMovementService.Service) *ScheduledMovementHandler {
	return &ScheduledMovementHandler{scheduledMovementService: scheduledMovementService}
}

func (h *ScheduledMovementHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	movement, ok := parseScheduledMovementRequest(c, "")
	if !ok {
		return nil
	}

	if err := h.scheduledMovementService.Create(c.Context(), userID, movement); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create scheduled movement")
	}
	return response.Created(c, mapScheduledMovementToResponse(movement))
}

func (h *ScheduledMovementHandler) GetAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	movements, err := h.scheduledMovementService.GetAll(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get scheduled movements")
	}

	resp := make([]dto.ScheduledMovementResponse, len(movements))
	for i, movement := range movements {
		resp[i] = mapScheduledMovementToResponse(movement)
	}
	return response.JSON(c, fiber.StatusOK, resp)
}

func (h *ScheduledMovementHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	movement, err := h.scheduledMovementService.GetByID(c.Context(), userID, c.Params("id"))
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "scheduled movement not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to get scheduled movement")
	}
	return response.JSON(c, fiber.StatusOK, mapScheduledMovementToResponse(movement))
}

func (h *ScheduledMovementHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	movement, ok := parseScheduledMovementRequest(c, c.Params("id"))
	if !ok {
		return nil
	}

	if err := h.scheduledMovementService.Update(c.Context(), userID, movement); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "scheduled movement not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to update scheduled movement")
	}
	return response.JSON(c, fiber.StatusOK, mapScheduledMovementToResponse(movement))
}

func (h *ScheduledMovementHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if err := h.scheduledMovementService.Delete(c.Context(), userID, c.Params("id")); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "scheduled movement not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete scheduled movement")
	}
	return response.JSON(c, fiber.StatusOK, fiber.Map{"message": "scheduled movement deleted successfully"})
}

func (h *ScheduledMovementHandler) GenerateDue(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	today := time.Now()
	if dateParam := c.Query("date"); dateParam != "" {
		parsed, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		}
		today = parsed
	}

	transactions, err := h.scheduledMovementService.GenerateDue(c.Context(), userID, today)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to generate scheduled movements")
	}

	items := make([]dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		items[i] = mapTransactionToResponse(tx)
	}

	return response.JSON(c, fiber.StatusOK, dto.GenerateScheduledMovementsResponse{
		Created: len(items),
		Items:   items,
	})
}

func parseScheduledMovementRequest(c *fiber.Ctx, id string) (*domain.ScheduledMovement, bool) {
	var req dto.CreateScheduledMovementRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			_ = response.ValidationError(c, vErr.Errors)
			return nil, false
		}
		_ = response.Error(c, fiber.StatusBadRequest, err.Error())
		return nil, false
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		_ = response.Error(c, fiber.StatusBadRequest, "invalid start_date format, use YYYY-MM-DD")
		return nil, false
	}

	nextRunDate := startDate
	if req.NextRunDate != nil && *req.NextRunDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.NextRunDate)
		if err != nil {
			_ = response.Error(c, fiber.StatusBadRequest, "invalid next_run_date format, use YYYY-MM-DD")
			return nil, false
		}
		nextRunDate = parsed
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			_ = response.Error(c, fiber.StatusBadRequest, "invalid end_date format, use YYYY-MM-DD")
			return nil, false
		}
		endDate = &parsed
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	return &domain.ScheduledMovement{
		ID:          id,
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Type:        domain.TransactionType(req.Type),
		Amount:      req.Amount,
		Description: req.Description,
		Notes:       req.Notes,
		Frequency:   domain.ScheduledMovementFrequency(req.Frequency),
		StartDate:   startDate,
		NextRunDate: nextRunDate,
		EndDate:     endDate,
		Active:      active,
	}, true
}

func mapScheduledMovementToResponse(movement *domain.ScheduledMovement) dto.ScheduledMovementResponse {
	var endDate *string
	if movement.EndDate != nil {
		s := movement.EndDate.Format("2006-01-02")
		endDate = &s
	}

	var lastGeneratedDate *string
	if movement.LastGeneratedDate != nil {
		s := movement.LastGeneratedDate.Format("2006-01-02")
		lastGeneratedDate = &s
	}

	return dto.ScheduledMovementResponse{
		ID:                movement.ID,
		UserID:            movement.UserID,
		AccountID:         movement.AccountID,
		CategoryID:        movement.CategoryID,
		Type:              string(movement.Type),
		Amount:            movement.Amount,
		Description:       movement.Description,
		Notes:             movement.Notes,
		Frequency:         string(movement.Frequency),
		StartDate:         movement.StartDate.Format("2006-01-02"),
		NextRunDate:       movement.NextRunDate.Format("2006-01-02"),
		EndDate:           endDate,
		Active:            movement.Active,
		LastGeneratedDate: lastGeneratedDate,
		CreatedAt:         movement.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         movement.UpdatedAt.Format(time.RFC3339),
	}
}
