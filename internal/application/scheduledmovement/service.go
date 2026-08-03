package scheduledmovement

import (
	"context"
	"fmt"
	"time"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	Create(ctx context.Context, userID string, movement *domain.ScheduledMovement) error
	GetByID(ctx context.Context, userID, id string) (*domain.ScheduledMovement, error)
	GetAll(ctx context.Context, userID string) ([]*domain.ScheduledMovement, error)
	Update(ctx context.Context, userID string, movement *domain.ScheduledMovement) error
	Delete(ctx context.Context, userID, id string) error
	GenerateDue(ctx context.Context, userID string, today time.Time) ([]*domain.Transaction, error)
}

type service struct {
	movementRepo    domain.ScheduledMovementRepository
	transactionRepo domain.TransactionRepository
}

func NewService(movementRepo domain.ScheduledMovementRepository, transactionRepo domain.TransactionRepository) Service {
	return &service{movementRepo: movementRepo, transactionRepo: transactionRepo}
}

func (s *service) Create(ctx context.Context, userID string, movement *domain.ScheduledMovement) error {
	movement.UserID = userID
	if movement.NextRunDate.IsZero() {
		movement.NextRunDate = movement.StartDate
	}
	if err := normalizeScheduledMovement(movement); err != nil {
		return err
	}
	return s.movementRepo.Create(ctx, movement)
}

func (s *service) GetByID(ctx context.Context, userID, id string) (*domain.ScheduledMovement, error) {
	movement, err := s.movementRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if movement.UserID != userID {
		return nil, domain.ErrNotFound
	}
	return movement, nil
}

func (s *service) GetAll(ctx context.Context, userID string) ([]*domain.ScheduledMovement, error) {
	return s.movementRepo.GetByUserID(ctx, userID)
}

func (s *service) Update(ctx context.Context, userID string, movement *domain.ScheduledMovement) error {
	existing, err := s.GetByID(ctx, userID, movement.ID)
	if err != nil {
		return err
	}
	movement.UserID = userID
	if movement.NextRunDate.IsZero() {
		movement.NextRunDate = existing.NextRunDate
	}
	if err := normalizeScheduledMovement(movement); err != nil {
		return err
	}
	return s.movementRepo.Update(ctx, movement)
}

func (s *service) Delete(ctx context.Context, userID, id string) error {
	if _, err := s.GetByID(ctx, userID, id); err != nil {
		return err
	}
	return s.movementRepo.Delete(ctx, id)
}

func (s *service) GenerateDue(ctx context.Context, userID string, today time.Time) ([]*domain.Transaction, error) {
	due, err := s.movementRepo.GetDueByUserID(ctx, userID, today)
	if err != nil {
		return nil, err
	}

	created := make([]*domain.Transaction, 0)
	for _, movement := range due {
		runDate := movement.NextRunDate
		for !runDate.After(today) && (movement.EndDate == nil || !runDate.After(*movement.EndDate)) {
			tx := &domain.Transaction{
				UserID:      userID,
				AccountID:   movement.AccountID,
				CategoryID:  movement.CategoryID,
				Type:        movement.Type,
				Amount:      movement.Amount,
				Description: movement.Description,
				Notes:       movement.Notes,
				Date:        runDate,
			}
			if err := s.transactionRepo.Create(ctx, tx); err != nil {
				return nil, fmt.Errorf("failed to create scheduled transaction: %w", err)
			}
			created = append(created, tx)

			lastRun := runDate
			movement.LastGeneratedDate = &lastRun
			runDate = nextRunDate(runDate, movement.Frequency)
		}

		movement.NextRunDate = runDate
		if movement.EndDate != nil && runDate.After(*movement.EndDate) {
			movement.Active = false
		}
		if err := s.movementRepo.Update(ctx, movement); err != nil {
			return nil, fmt.Errorf("failed to update scheduled movement: %w", err)
		}
	}

	return created, nil
}

func nextRunDate(date time.Time, frequency domain.ScheduledMovementFrequency) time.Time {
	switch frequency {
	case domain.ScheduledMovementWeekdays:
		next := date.AddDate(0, 0, 1)
		for next.Weekday() == time.Saturday || next.Weekday() == time.Sunday {
			next = next.AddDate(0, 0, 1)
		}
		return next
	case domain.ScheduledMovementWeekly:
		return date.AddDate(0, 0, 7)
	case domain.ScheduledMovementMonthly:
		return date.AddDate(0, 1, 0)
	case domain.ScheduledMovementYearly:
		return date.AddDate(1, 0, 0)
	default:
		return date.AddDate(0, 0, 1)
	}
}

func normalizeScheduledMovement(movement *domain.ScheduledMovement) error {
	if movement.Type == domain.TransactionTypeInformational {
		movement.Amount = 0
		movement.CategoryID = nil
		return nil
	}

	if movement.Amount <= 0 {
		return fmt.Errorf("%w: amount must be greater than zero", domain.ErrValidation)
	}
	return nil
}
