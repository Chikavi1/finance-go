package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	reminderService "github.com/agnathor/finances-go/internal/application/reminder"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/internal/interfaces/api/validator"
	"github.com/agnathor/finances-go/pkg/response"
)

type ReminderHandler struct {
	reminderService reminderService.Service
}

func NewReminderHandler(reminderService reminderService.Service) *ReminderHandler {
	return &ReminderHandler{reminderService: reminderService}
}

func (h *ReminderHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req dto.CreateReminderRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid due_date format, use YYYY-MM-DD")
	}
	reminderTime, err := parseReminderTime(req.ReminderTime)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid reminder_time format, use HH:MM")
	}

	dayOfMonth := resolveReminderDayOfMonth(req.DayOfMonth, dueDate)
	reminder, err := h.reminderService.Create(c.Context(), userID, req.Title, req.Amount, dueDate, reminderTime, domain.ReminderRecurrenceType(req.RecurrenceType), dayOfMonth, domain.ReminderStatus(req.Status), req.RelatedType, req.RelatedID, req.Notes)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create reminder")
	}

	return response.Created(c, mapReminderToResponse(reminder))
}

func (h *ReminderHandler) GetAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	reminders, err := h.reminderService.GetAll(c.Context(), userID, c.QueryBool("include_done"))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get reminders")
	}

	resp := make([]dto.ReminderResponse, len(reminders))
	for i, reminder := range reminders {
		resp[i] = mapReminderToResponse(reminder)
	}
	return response.JSON(c, fiber.StatusOK, resp)
}

func (h *ReminderHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	reminder, err := h.reminderService.GetByID(c.Context(), userID, c.Params("id"))
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "reminder not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to get reminder")
	}
	return response.JSON(c, fiber.StatusOK, mapReminderToResponse(reminder))
}

func (h *ReminderHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var req dto.UpdateReminderRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid due_date format, use YYYY-MM-DD")
	}
	reminderTime, err := parseReminderTime(req.ReminderTime)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid reminder_time format, use HH:MM")
	}

	dayOfMonth := resolveReminderDayOfMonth(req.DayOfMonth, dueDate)
	reminder, err := h.reminderService.Update(c.Context(), userID, c.Params("id"), req.Title, req.Amount, dueDate, reminderTime, domain.ReminderRecurrenceType(req.RecurrenceType), dayOfMonth, domain.ReminderStatus(req.Status), req.RelatedType, req.RelatedID, req.Notes)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "reminder not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to update reminder")
	}
	return response.JSON(c, fiber.StatusOK, mapReminderToResponse(reminder))
}

func (h *ReminderHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if err := h.reminderService.Delete(c.Context(), userID, c.Params("id")); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "reminder not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete reminder")
	}
	return response.JSON(c, fiber.StatusOK, fiber.Map{"message": "reminder deleted successfully"})
}

func mapReminderToResponse(reminder *domain.Reminder) dto.ReminderResponse {
	var notificationSentAt *string
	if reminder.NotificationSentAt != nil {
		s := reminder.NotificationSentAt.Format(time.RFC3339)
		notificationSentAt = &s
	}

	return dto.ReminderResponse{
		ID:          reminder.ID,
		UserID:      reminder.UserID,
		Title:       reminder.Title,
		Amount:      reminder.Amount,
		DueDate:     reminder.DueDate.Format("2006-01-02"),
		ReminderTime: reminder.ReminderTime,
		RecurrenceType: string(reminder.RecurrenceType),
		DayOfMonth: reminder.DayOfMonth,
		Status:      string(reminder.Status),
		RelatedType: reminder.RelatedType,
		RelatedID:   reminder.RelatedID,
		Notes:       reminder.Notes,
		NotificationSentAt: notificationSentAt,
		CreatedAt:   reminder.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   reminder.UpdatedAt.Format(time.RFC3339),
	}
}

func resolveReminderDayOfMonth(value *int, dueDate time.Time) *int {
	if value != nil {
		return value
	}
	day := dueDate.Day()
	return &day
}

func parseReminderTime(value string) (string, error) {
	if value == "" {
		return "09:00", nil
	}
	parsed, err := time.Parse("15:04", value)
	if err != nil {
		return "", err
	}
	return parsed.Format("15:04"), nil
}
