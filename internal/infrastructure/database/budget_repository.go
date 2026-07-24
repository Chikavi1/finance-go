package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
)

type budgetRepository struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func NewBudgetRepository(pool *pgxpool.Pool) domain.BudgetRepository {
	return &budgetRepository{
		pool:  pool,
		query: db.New(pool),
	}
}

func (r *budgetRepository) Create(ctx context.Context, budget *domain.Budget) error {
	userUUID, err := parseUUID(budget.UserID)
	if err != nil {
		return domain.ErrNotFound
	}

	categoryUUID, err := parseUUID(budget.CategoryID)
	if err != nil {
		return domain.ErrNotFound
	}

	created, err := r.query.CreateBudget(ctx, db.CreateBudgetParams{
		UserID:     userUUID,
		CategoryID: categoryUUID,
		Amount:     budget.Amount,
		Spent:      budget.Spent,
		Month:      budget.Month,
		Year:       budget.Year,
	})
	if err != nil {
		return err
	}

	budget.ID = pgUUIDToString(created.ID)
	budget.CreatedAt = created.CreatedAt.Time
	budget.UpdatedAt = created.UpdatedAt.Time
	return nil
}

func (r *budgetRepository) GetByID(ctx context.Context, id string) (*domain.Budget, error) {
	budgetUUID, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	budget, err := r.query.GetBudgetByID(ctx, budgetUUID)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return mapBudget(budget), nil
}

func (r *budgetRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Budget, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	budgets, err := r.query.GetBudgetsByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Budget, len(budgets))
	for i, b := range budgets {
		result[i] = mapBudget(b)
	}

	return result, nil
}

func (r *budgetRepository) GetByMonthYear(ctx context.Context, userID string, month, year int32) ([]*domain.Budget, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	budgets, err := r.query.GetBudgetsByMonthYear(ctx, db.GetBudgetsByMonthYearParams{
		UserID: userUUID,
		Month:  month,
		Year:   year,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Budget, len(budgets))
	for i, b := range budgets {
		result[i] = mapBudget(b)
	}

	return result, nil
}

func (r *budgetRepository) Update(ctx context.Context, budget *domain.Budget) error {
	budgetUUID, err := parseUUID(budget.ID)
	if err != nil {
		return domain.ErrNotFound
	}

	updated, err := r.query.UpdateBudget(ctx, db.UpdateBudgetParams{
		ID:     budgetUUID,
		Amount: budget.Amount,
		Spent:  budget.Spent,
	})
	if err != nil {
		if isNotFound(err) {
			return domain.ErrNotFound
		}
		return err
	}

	budget.UpdatedAt = updated.UpdatedAt.Time
	return nil
}

func (r *budgetRepository) Delete(ctx context.Context, id string) error {
	budgetUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	return r.query.DeleteBudget(ctx, budgetUUID)
}
