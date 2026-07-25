package domain

import (
	"context"
	"time"
)

type Goal struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	Name          string     `json:"name"`
	TargetAmount  float64    `json:"target_amount"`
	CurrentAmount float64    `json:"current_amount"`
	TargetDate    *time.Time `json:"target_date,omitempty"`
	Icon          string     `json:"icon"`
	Color         string     `json:"color"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type GoalContribution struct {
	ID               string    `json:"id"`
	GoalID           string    `json:"goal_id"`
	Amount           float64   `json:"amount"`
	ContributionDate time.Time `json:"contribution_date"`
	Notes            *string   `json:"notes,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type GoalRepository interface {
	Create(ctx context.Context, goal *Goal) error
	GetByID(ctx context.Context, id string) (*Goal, error)
	GetByUserID(ctx context.Context, userID string) ([]*Goal, error)
	Update(ctx context.Context, goal *Goal) error
	Delete(ctx context.Context, id string) error
}

type GoalContributionRepository interface {
	Create(ctx context.Context, contribution *GoalContribution) error
	GetByGoalID(ctx context.Context, goalID string) ([]*GoalContribution, error)
	Delete(ctx context.Context, id string) error
}
