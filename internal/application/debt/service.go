package debt

import (
	"context"
	"fmt"
	"time"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	Create(ctx context.Context, userID, name string, totalAmount, remainingAmount, interestRate float64, dueDate *time.Time, status domain.DebtStatus, notes *string) (*domain.Debt, error)
	GetByID(ctx context.Context, userID, debtID string) (*domain.Debt, error)
	GetAll(ctx context.Context, userID string) ([]*domain.Debt, error)
	Update(ctx context.Context, userID, debtID, name string, totalAmount, remainingAmount, interestRate float64, dueDate *time.Time, status domain.DebtStatus, notes *string) (*domain.Debt, error)
	Delete(ctx context.Context, userID, debtID string) error
}

type service struct {
	debtRepo domain.DebtRepository
}

func NewService(debtRepo domain.DebtRepository) Service {
	return &service{debtRepo: debtRepo}
}

func (s *service) Create(ctx context.Context, userID, name string, totalAmount, remainingAmount, interestRate float64, dueDate *time.Time, status domain.DebtStatus, notes *string) (*domain.Debt, error) {
	debt := &domain.Debt{
		UserID:          userID,
		Name:            name,
		TotalAmount:     totalAmount,
		RemainingAmount: remainingAmount,
		InterestRate:    interestRate,
		DueDate:         dueDate,
		Status:          status,
		Notes:           notes,
	}

	if err := s.debtRepo.Create(ctx, debt); err != nil {
		return nil, fmt.Errorf("failed to create debt: %w", err)
	}

	return debt, nil
}

func (s *service) GetByID(ctx context.Context, userID, debtID string) (*domain.Debt, error) {
	debt, err := s.debtRepo.GetByID(ctx, debtID)
	if err != nil {
		return nil, err
	}

	if debt.UserID != userID {
		return nil, domain.ErrNotFound
	}

	return debt, nil
}

func (s *service) GetAll(ctx context.Context, userID string) ([]*domain.Debt, error) {
	return s.debtRepo.GetByUserID(ctx, userID)
}

func (s *service) Update(ctx context.Context, userID, debtID, name string, totalAmount, remainingAmount, interestRate float64, dueDate *time.Time, status domain.DebtStatus, notes *string) (*domain.Debt, error) {
	debt, err := s.debtRepo.GetByID(ctx, debtID)
	if err != nil {
		return nil, err
	}

	if debt.UserID != userID {
		return nil, domain.ErrNotFound
	}

	debt.Name = name
	debt.TotalAmount = totalAmount
	debt.RemainingAmount = remainingAmount
	debt.InterestRate = interestRate
	debt.DueDate = dueDate
	debt.Status = status
	debt.Notes = notes

	if err := s.debtRepo.Update(ctx, debt); err != nil {
		return nil, fmt.Errorf("failed to update debt: %w", err)
	}

	return debt, nil
}

func (s *service) Delete(ctx context.Context, userID, debtID string) error {
	debt, err := s.debtRepo.GetByID(ctx, debtID)
	if err != nil {
		return err
	}

	if debt.UserID != userID {
		return domain.ErrNotFound
	}

	return s.debtRepo.Delete(ctx, debtID)
}
