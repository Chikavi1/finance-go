package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/agnathor/finances-go/internal/application/goal"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/internal/interfaces/api/dto"
	"github.com/agnathor/finances-go/internal/interfaces/api/middleware"
	"github.com/agnathor/finances-go/internal/interfaces/api/validator"
	"github.com/agnathor/finances-go/pkg/response"
)

type GoalHandler struct {
	goalService goal.Service
}

func NewGoalHandler(goalService goal.Service) *GoalHandler {
	return &GoalHandler{goalService: goalService}
}

func (h *GoalHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req dto.CreateGoalRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	var targetDate *time.Time
	if req.TargetDate != nil {
		t, err := time.Parse("2006-01-02", *req.TargetDate)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid target_date format, use YYYY-MM-DD")
		}
		targetDate = &t
	}

	goal, err := h.goalService.Create(c.Context(), userID, req.Name, req.TargetAmount, req.CurrentAmount, targetDate, req.Icon, req.Color)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create goal")
	}

	return response.Created(c, mapGoalToResponse(goal))
}

func (h *GoalHandler) GetByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	goalID := c.Params("id")

	goal, err := h.goalService.GetByID(c.Context(), userID, goalID)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "goal not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to get goal")
	}

	return response.JSON(c, fiber.StatusOK, mapGoalToResponse(goal))
}

func (h *GoalHandler) GetAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	goals, err := h.goalService.GetAll(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get goals")
	}

	resp := make([]dto.GoalResponse, len(goals))
	for i, g := range goals {
		resp[i] = mapGoalToResponse(g)
	}

	return response.JSON(c, fiber.StatusOK, resp)
}

func (h *GoalHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	goalID := c.Params("id")

	var req dto.UpdateGoalRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, vErr.Errors)
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	var targetDate *time.Time
	if req.TargetDate != nil {
		t, err := time.Parse("2006-01-02", *req.TargetDate)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid target_date format, use YYYY-MM-DD")
		}
		targetDate = &t
	}

	goal, err := h.goalService.Update(c.Context(), userID, goalID, req.Name, req.TargetAmount, req.CurrentAmount, targetDate, req.Icon, req.Color)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "goal not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to update goal")
	}

	return response.JSON(c, fiber.StatusOK, mapGoalToResponse(goal))
}

func (h *GoalHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	goalID := c.Params("id")

	if err := h.goalService.Delete(c.Context(), userID, goalID); err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "goal not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete goal")
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"message": "goal deleted successfully",
	})
}

func mapGoalToResponse(g *domain.Goal) dto.GoalResponse {
	var targetDate *string
	if g.TargetDate != nil {
		t := g.TargetDate.Format("2006-01-02")
		targetDate = &t
	}

	return dto.GoalResponse{
		ID:            g.ID,
		UserID:        g.UserID,
		Name:          g.Name,
		TargetAmount:  g.TargetAmount,
		CurrentAmount: g.CurrentAmount,
		TargetDate:    targetDate,
		Icon:          g.Icon,
		Color:         g.Color,
		CreatedAt:     g.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     g.UpdatedAt.Format(time.RFC3339),
	}
}
