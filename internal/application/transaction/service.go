package transaction

import (
	"context"

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
