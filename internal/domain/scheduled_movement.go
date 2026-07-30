package domain

import (
	"context"
	"time"
)

type ScheduledMovementFrequency string

const (
	ScheduledMovementDaily    ScheduledMovementFrequency = "daily"
	ScheduledMovementWeekdays ScheduledMovementFrequency = "weekdays"
	ScheduledMovementWeekly   ScheduledMovementFrequency = "weekly"
	ScheduledMovementMonthly  ScheduledMovementFrequency = "monthly"
	ScheduledMovementYearly   ScheduledMovementFrequency = "yearly"
)

type ScheduledMovement struct {
	ID                string                     `json:"id"`
	UserID            string                     `json:"user_id"`
	AccountID         string                     `json:"account_id"`
	CategoryID        *string                    `json:"category_id,omitempty"`
	Type              TransactionType            `json:"type"`
	Amount            float64                    `json:"amount"`
	Description       string                     `json:"description"`
	Notes             *string                    `json:"notes,omitempty"`
	Frequency         ScheduledMovementFrequency `json:"frequency"`
	StartDate         time.Time                  `json:"start_date"`
	NextRunDate       time.Time                  `json:"next_run_date"`
	EndDate           *time.Time                 `json:"end_date,omitempty"`
	Active            bool                       `json:"active"`
	LastGeneratedDate *time.Time                 `json:"last_generated_date,omitempty"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
}

type ScheduledMovementRepository interface {
	Create(ctx context.Context, movement *ScheduledMovement) error
	GetByID(ctx context.Context, id string) (*ScheduledMovement, error)
	GetByUserID(ctx context.Context, userID string) ([]*ScheduledMovement, error)
	GetDueByUserID(ctx context.Context, userID string, dueDate time.Time) ([]*ScheduledMovement, error)
	Update(ctx context.Context, movement *ScheduledMovement) error
	Delete(ctx context.Context, id string) error
}
