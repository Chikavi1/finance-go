package account

import (
	"context"
	"fmt"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	Create(ctx context.Context, userID string, name string, accountType domain.AccountType, currency string, balance float64, color string, icon string) (*domain.Account, error)
	GetByID(ctx context.Context, userID, accountID string) (*domain.Account, error)
	GetAll(ctx context.Context, userID string) ([]*domain.Account, error)
	Update(ctx context.Context, userID string, accountID string, name string, accountType domain.AccountType, currency string, color string, icon string, archived bool) (*domain.Account, error)
	Delete(ctx context.Context, userID, accountID string) error
}

type service struct {
	accountRepo domain.AccountRepository
}

func NewService(accountRepo domain.AccountRepository) Service {
	return &service{accountRepo: accountRepo}
}

func (s *service) Create(ctx context.Context, userID string, name string, accountType domain.AccountType, currency string, balance float64, color string, icon string) (*domain.Account, error) {
	account := &domain.Account{
		UserID:   userID,
		Name:     name,
		Type:     accountType,
		Currency: currency,
		Balance:  balance,
		Color:    color,
		Icon:     icon,
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (s *service) GetByID(ctx context.Context, userID, accountID string) (*domain.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if account.UserID != userID {
		return nil, domain.ErrNotFound
	}

	return account, nil
}

func (s *service) GetAll(ctx context.Context, userID string) ([]*domain.Account, error) {
	return s.accountRepo.GetByUserID(ctx, userID)
}

func (s *service) Update(ctx context.Context, userID string, accountID string, name string, accountType domain.AccountType, currency string, color string, icon string, archived bool) (*domain.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if account.UserID != userID {
		return nil, domain.ErrNotFound
	}

	account.Name = name
	account.Type = accountType
	account.Currency = currency
	account.Color = color
	account.Icon = icon
	account.Archived = archived

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return account, nil
}

func (s *service) Delete(ctx context.Context, userID, accountID string) error {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}

	if account.UserID != userID {
		return domain.ErrNotFound
	}

	return s.accountRepo.Delete(ctx, accountID)
}
