package domain

import (
	"context"
	"time"
)

type DebtStatus string

const (
	DebtStatusActive  DebtStatus = "active"
	DebtStatusPaid    DebtStatus = "paid"
	DebtStatusOverdue DebtStatus = "overdue"
)

type Debt struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	Name            string     `json:"name"`
	TotalAmount     float64    `json:"total_amount"`
	RemainingAmount float64    `json:"remaining_amount"`
	InterestRate    float64    `json:"interest_rate"`
	DueDate         *time.Time `json:"due_date,omitempty"`
	Status          DebtStatus `json:"status"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type DebtPayment struct {
	ID          string    `json:"id"`
	DebtID      string    `json:"debt_id"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
	Notes       *string   `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type DebtRepository interface {
	Create(ctx context.Context, debt *Debt) error
	GetByID(ctx context.Context, id string) (*Debt, error)
	GetByUserID(ctx context.Context, userID string) ([]*Debt, error)
	Update(ctx context.Context, debt *Debt) error
	Delete(ctx context.Context, id string) error
}

type DebtPaymentRepository interface {
	Create(ctx context.Context, payment *DebtPayment) error
	GetByDebtID(ctx context.Context, debtID string) ([]*DebtPayment, error)
	Delete(ctx context.Context, id string) error
}
