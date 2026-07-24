package goal

import (
	"context"
	"fmt"
	"time"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	Create(ctx context.Context, userID, name string, targetAmount, currentAmount float64, targetDate *time.Time, icon, color string) (*domain.Goal, error)
	GetByID(ctx context.Context, userID, goalID string) (*domain.Goal, error)
	GetAll(ctx context.Context, userID string) ([]*domain.Goal, error)
	Update(ctx context.Context, userID, goalID, name string, targetAmount, currentAmount float64, targetDate *time.Time, icon, color string) (*domain.Goal, error)
	Delete(ctx context.Context, userID, goalID string) error
}

type service struct {
	goalRepo domain.GoalRepository
}

func NewService(goalRepo domain.GoalRepository) Service {
	return &service{goalRepo: goalRepo}
}

func (s *service) Create(ctx context.Context, userID, name string, targetAmount, currentAmount float64, targetDate *time.Time, icon, color string) (*domain.Goal, error) {
	goal := &domain.Goal{
		UserID:        userID,
		Name:          name,
		TargetAmount:  targetAmount,
		CurrentAmount: currentAmount,
		TargetDate:    targetDate,
		Icon:          icon,
		Color:         color,
	}

	if err := s.goalRepo.Create(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}

	return goal, nil
}

func (s *service) GetByID(ctx context.Context, userID, goalID string) (*domain.Goal, error) {
	goal, err := s.goalRepo.GetByID(ctx, goalID)
	if err != nil {
		return nil, err
	}

	if goal.UserID != userID {
		return nil, domain.ErrNotFound
	}

	return goal, nil
}

func (s *service) GetAll(ctx context.Context, userID string) ([]*domain.Goal, error) {
	return s.goalRepo.GetByUserID(ctx, userID)
}

func (s *service) Update(ctx context.Context, userID, goalID, name string, targetAmount, currentAmount float64, targetDate *time.Time, icon, color string) (*domain.Goal, error) {
	goal, err := s.goalRepo.GetByID(ctx, goalID)
	if err != nil {
		return nil, err
	}

	if goal.UserID != userID {
		return nil, domain.ErrNotFound
	}

	goal.Name = name
	goal.TargetAmount = targetAmount
	goal.CurrentAmount = currentAmount
	goal.TargetDate = targetDate
	goal.Icon = icon
	goal.Color = color

	if err := s.goalRepo.Update(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	return goal, nil
}

func (s *service) Delete(ctx context.Context, userID, goalID string) error {
	goal, err := s.goalRepo.GetByID(ctx, goalID)
	if err != nil {
		return err
	}

	if goal.UserID != userID {
		return domain.ErrNotFound
	}

	return s.goalRepo.Delete(ctx, goalID)
}
