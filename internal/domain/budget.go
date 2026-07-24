package domain

import (
	"context"
	"time"
)

type Budget struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CategoryID string    `json:"category_id"`
	Amount     float64   `json:"amount"`
	Spent      float64   `json:"spent"`
	Month      int32     `json:"month"`
	Year       int32     `json:"year"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BudgetRepository interface {
	Create(ctx context.Context, budget *Budget) error
	GetByID(ctx context.Context, id string) (*Budget, error)
	GetByUserID(ctx context.Context, userID string) ([]*Budget, error)
	GetByMonthYear(ctx context.Context, userID string, month, year int32) ([]*Budget, error)
	Update(ctx context.Context, budget *Budget) error
	Delete(ctx context.Context, id string) error
}
