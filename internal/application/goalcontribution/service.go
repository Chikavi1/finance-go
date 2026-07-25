package goalcontribution

import (
	"context"
	"fmt"
	"time"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	Create(ctx context.Context, goalID string, amount float64, contributionDate time.Time, notes *string) (*domain.GoalContribution, error)
	GetByGoalID(ctx context.Context, goalID string) ([]*domain.GoalContribution, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	goalContributionRepo domain.GoalContributionRepository
	goalRepo             domain.GoalRepository
}

func NewService(goalContributionRepo domain.GoalContributionRepository, goalRepo domain.GoalRepository) Service {
	return &service{
		goalContributionRepo: goalContributionRepo,
		goalRepo:             goalRepo,
	}
}

func (s *service) Create(ctx context.Context, goalID string, amount float64, contributionDate time.Time, notes *string) (*domain.GoalContribution, error) {
	contribution := &domain.GoalContribution{
		GoalID:           goalID,
		Amount:           amount,
		ContributionDate: contributionDate,
		Notes:            notes,
	}

	if err := s.goalContributionRepo.Create(ctx, contribution); err != nil {
		return nil, fmt.Errorf("failed to create goal contribution: %w", err)
	}

	if goalObj, err := s.goalRepo.GetByID(ctx, goalID); err == nil && goalObj != nil {
		goalObj.CurrentAmount += amount
		_ = s.goalRepo.Update(ctx, goalObj)
	}

	return contribution, nil
}

func (s *service) GetByGoalID(ctx context.Context, goalID string) ([]*domain.GoalContribution, error) {
	return s.goalContributionRepo.GetByGoalID(ctx, goalID)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.goalContributionRepo.Delete(ctx, id)
}
