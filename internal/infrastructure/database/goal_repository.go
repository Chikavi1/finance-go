package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
)

type goalRepository struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func NewGoalRepository(pool *pgxpool.Pool) domain.GoalRepository {
	return &goalRepository{
		pool:  pool,
		query: db.New(pool),
	}
}

func (r *goalRepository) Create(ctx context.Context, goal *domain.Goal) error {
	userUUID, err := parseUUID(goal.UserID)
	if err != nil {
		return domain.ErrNotFound
	}

	created, err := r.query.CreateGoal(ctx, db.CreateGoalParams{
		UserID:        userUUID,
		Name:          goal.Name,
		TargetAmount:  goal.TargetAmount,
		CurrentAmount: goal.CurrentAmount,
		TargetDate:    toNullableDate(goal.TargetDate),
		Icon:          goal.Icon,
		Color:         goal.Color,
	})
	if err != nil {
		return err
	}

	goal.ID = pgUUIDToString(created.ID)
	goal.CreatedAt = created.CreatedAt.Time
	goal.UpdatedAt = created.UpdatedAt.Time
	return nil
}

func (r *goalRepository) GetByID(ctx context.Context, id string) (*domain.Goal, error) {
	goalUUID, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	goal, err := r.query.GetGoalByID(ctx, goalUUID)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return mapGoal(goal), nil
}

func (r *goalRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Goal, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	goals, err := r.query.GetGoalsByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Goal, len(goals))
	for i, g := range goals {
		result[i] = mapGoal(g)
	}

	return result, nil
}

func (r *goalRepository) Update(ctx context.Context, goal *domain.Goal) error {
	goalUUID, err := parseUUID(goal.ID)
	if err != nil {
		return domain.ErrNotFound
	}

	updated, err := r.query.UpdateGoal(ctx, db.UpdateGoalParams{
		ID:            goalUUID,
		Name:          goal.Name,
		TargetAmount:  goal.TargetAmount,
		CurrentAmount: goal.CurrentAmount,
		TargetDate:    toNullableDate(goal.TargetDate),
		Icon:          goal.Icon,
		Color:         goal.Color,
	})
	if err != nil {
		if isNotFound(err) {
			return domain.ErrNotFound
		}
		return err
	}

	goal.UpdatedAt = updated.UpdatedAt.Time
	return nil
}

func (r *goalRepository) Delete(ctx context.Context, id string) error {
	goalUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	return r.query.DeleteGoal(ctx, goalUUID)
}
