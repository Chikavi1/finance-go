package domain

import (
	"context"
	"time"
)

type AccountType string

const (
	AccountTypeCash       AccountType = "cash"
	AccountTypeWallet     AccountType = "wallet"
	AccountTypeBank       AccountType = "bank"
	AccountTypeCreditCard AccountType = "credit_card"
	AccountTypeSavings    AccountType = "savings"
	AccountTypeInvestment AccountType = "investment"
)

type Account struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Name      string      `json:"name"`
	Type      AccountType `json:"type"`
	Currency  string      `json:"currency"`
	Balance   float64     `json:"balance"`
	Color     string      `json:"color"`
	Icon      string      `json:"icon"`
	Archived  bool        `json:"archived"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	GetByID(ctx context.Context, id string) (*Account, error)
	GetByUserID(ctx context.Context, userID string) ([]*Account, error)
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, id string) error
}
