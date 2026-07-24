package budget

import (
	"context"
	"fmt"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	Create(ctx context.Context, userID, categoryID string, amount float64, month, year int32) (*domain.Budget, error)
	GetByID(ctx context.Context, userID, budgetID string) (*domain.Budget, error)
	GetAll(ctx context.Context, userID string) ([]*domain.Budget, error)
	GetByMonthYear(ctx context.Context, userID string, month, year int32) ([]*domain.Budget, error)
	Update(ctx context.Context, userID, budgetID string, amount, spent float64) (*domain.Budget, error)
	Delete(ctx context.Context, userID, budgetID string) error
}

type service struct {
	budgetRepo domain.BudgetRepository
}

func NewService(budgetRepo domain.BudgetRepository) Service {
	return &service{budgetRepo: budgetRepo}
}

func (s *service) Create(ctx context.Context, userID, categoryID string, amount float64, month, year int32) (*domain.Budget, error) {
	budget := &domain.Budget{
		UserID:     userID,
		CategoryID: categoryID,
		Amount:     amount,
		Month:      month,
		Year:       year,
	}

	if err := s.budgetRepo.Create(ctx, budget); err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	return budget, nil
}

func (s *service) GetByID(ctx context.Context, userID, budgetID string) (*domain.Budget, error) {
	budget, err := s.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		return nil, err
	}

	if budget.UserID != userID {
		return nil, domain.ErrNotFound
	}

	return budget, nil
}

func (s *service) GetAll(ctx context.Context, userID string) ([]*domain.Budget, error) {
	return s.budgetRepo.GetByUserID(ctx, userID)
}

func (s *service) GetByMonthYear(ctx context.Context, userID string, month, year int32) ([]*domain.Budget, error) {
	return s.budgetRepo.GetByMonthYear(ctx, userID, month, year)
}

func (s *service) Update(ctx context.Context, userID, budgetID string, amount, spent float64) (*domain.Budget, error) {
	budget, err := s.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		return nil, err
	}

	if budget.UserID != userID {
		return nil, domain.ErrNotFound
	}

	budget.Amount = amount
	budget.Spent = spent

	if err := s.budgetRepo.Update(ctx, budget); err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}

	return budget, nil
}

func (s *service) Delete(ctx context.Context, userID, budgetID string) error {
	budget, err := s.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		return err
	}

	if budget.UserID != userID {
		return domain.ErrNotFound
	}

	return s.budgetRepo.Delete(ctx, budgetID)
}
