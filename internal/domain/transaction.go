package domain

import (
	"context"
	"time"
)

type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

type Transaction struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	AccountID   string          `json:"account_id"`
	ToAccountID *string         `json:"to_account_id,omitempty"`
	CategoryID  *string         `json:"category_id,omitempty"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	Description string          `json:"description"`
	Notes       *string         `json:"notes,omitempty"`
	Date        time.Time       `json:"date"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Tags        []string        `json:"tags,omitempty"`
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) error
	GetByID(ctx context.Context, id string) (*Transaction, error)
	GetByUserID(ctx context.Context, userID string, filter TransactionFilter) ([]*Transaction, error)
	Update(ctx context.Context, tx *Transaction) error
	Delete(ctx context.Context, id string) error
}

type TransactionFilter struct {
	Type      *TransactionType
	AccountID *string
	CategoryID *string
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	PerPage   int
}
