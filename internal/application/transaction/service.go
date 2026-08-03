package transaction

import (
	"context"
	"fmt"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	Create(ctx context.Context, userID string, tx *domain.Transaction) error
	GetByID(ctx context.Context, userID, id string) (*domain.Transaction, error)
	GetAll(ctx context.Context, userID string, filter domain.TransactionFilter) ([]*domain.Transaction, error)
	Update(ctx context.Context, userID string, tx *domain.Transaction) error
	Delete(ctx context.Context, userID, id string) error
}

type service struct {
	transactionRepo domain.TransactionRepository
}

func NewService(transactionRepo domain.TransactionRepository) Service {
	return &service{transactionRepo: transactionRepo}
}

func (s *service) Create(ctx context.Context, userID string, tx *domain.Transaction) error {
	tx.UserID = userID
	if err := normalizeTransaction(tx); err != nil {
		return err
	}
	return s.transactionRepo.Create(ctx, tx)
}

func (s *service) GetByID(ctx context.Context, userID, id string) (*domain.Transaction, error) {
	tx, err := s.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tx.UserID != userID {
		return nil, domain.ErrNotFound
	}
	return tx, nil
}

func (s *service) GetAll(ctx context.Context, userID string, filter domain.TransactionFilter) ([]*domain.Transaction, error) {
	return s.transactionRepo.GetByUserID(ctx, userID, filter)
}

func (s *service) Update(ctx context.Context, userID string, tx *domain.Transaction) error {
	existing, err := s.transactionRepo.GetByID(ctx, tx.ID)
	if err != nil {
		return err
	}
	if existing.UserID != userID {
		return domain.ErrNotFound
	}
	tx.UserID = userID
	if err := normalizeTransaction(tx); err != nil {
		return err
	}
	return s.transactionRepo.Update(ctx, tx)
}

func (s *service) Delete(ctx context.Context, userID, id string) error {
	existing, err := s.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.UserID != userID {
		return domain.ErrNotFound
	}
	return s.transactionRepo.Delete(ctx, id)
}

func normalizeTransaction(tx *domain.Transaction) error {
	if tx.Type == domain.TransactionTypeInformational {
		tx.Amount = 0
		tx.CategoryID = nil
		tx.ToAccountID = nil
		return nil
	}

	if tx.Amount <= 0 {
		return fmt.Errorf("%w: amount must be greater than zero", domain.ErrValidation)
	}
	return nil
}
