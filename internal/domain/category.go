package domain

import (
	"context"
	"time"
)

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

type Category struct {
	ID        string       `json:"id"`
	UserID    string       `json:"user_id"`
	Name      string       `json:"name"`
	Type      CategoryType `json:"type"`
	Color     string       `json:"color"`
	Icon      string       `json:"icon"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id string) (*Category, error)
	GetByUserID(ctx context.Context, userID string) ([]*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id string) error
}
