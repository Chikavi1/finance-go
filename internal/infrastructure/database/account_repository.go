package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
)

type accountRepository struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func NewAccountRepository(pool *pgxpool.Pool) domain.AccountRepository {
	return &accountRepository{
		pool:  pool,
		query: db.New(pool),
	}
}

func (r *accountRepository) Create(ctx context.Context, account *domain.Account) error {
	userUUID, err := parseUUID(account.UserID)
	if err != nil {
		return domain.ErrNotFound
	}

	created, err := r.query.CreateAccount(ctx, db.CreateAccountParams{
		UserID:   userUUID,
		Name:     account.Name,
		Type:     string(account.Type),
		Currency: account.Currency,
		Balance:  account.Balance,
		Color:    account.Color,
		Icon:     account.Icon,
	})
	if err != nil {
		return err
	}

	account.ID = pgUUIDToString(created.ID)
	account.CreatedAt = created.CreatedAt.Time
	account.UpdatedAt = created.UpdatedAt.Time
	return nil
}

func (r *accountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	accountUUID, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	account, err := r.query.GetAccountByID(ctx, accountUUID)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return mapAccount(account), nil
}

func (r *accountRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Account, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	accounts, err := r.query.GetAccountsByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Account, len(accounts))
	for i, a := range accounts {
		result[i] = mapAccount(a)
	}

	return result, nil
}

func (r *accountRepository) Update(ctx context.Context, account *domain.Account) error {
	accountUUID, err := parseUUID(account.ID)
	if err != nil {
		return domain.ErrNotFound
	}

	updated, err := r.query.UpdateAccount(ctx, db.UpdateAccountParams{
		ID:       accountUUID,
		Name:     account.Name,
		Type:     string(account.Type),
		Currency: account.Currency,
		Color:    account.Color,
		Icon:     account.Icon,
		Archived: account.Archived,
	})
	if err != nil {
		if isNotFound(err) {
			return domain.ErrNotFound
		}
		return err
	}

	account.Balance = updated.Balance
	account.UpdatedAt = updated.UpdatedAt.Time
	return nil
}

func (r *accountRepository) Delete(ctx context.Context, id string) error {
	accountUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	return r.query.DeleteAccount(ctx, accountUUID)
}
