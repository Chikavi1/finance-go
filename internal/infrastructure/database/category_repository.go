package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
)

type categoryRepository struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func NewCategoryRepository(pool *pgxpool.Pool) domain.CategoryRepository {
	return &categoryRepository{
		pool:  pool,
		query: db.New(pool),
	}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	userUUID, err := parseUUID(category.UserID)
	if err != nil {
		return domain.ErrNotFound
	}

	created, err := r.query.CreateCategory(ctx, db.CreateCategoryParams{
		UserID: userUUID,
		Name:   category.Name,
		Type:   string(category.Type),
		Color:  category.Color,
		Icon:   category.Icon,
	})
	if err != nil {
		return err
	}

	category.ID = pgUUIDToString(created.ID)
	category.CreatedAt = created.CreatedAt.Time
	category.UpdatedAt = created.UpdatedAt.Time
	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	categoryUUID, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	category, err := r.query.GetCategoryByID(ctx, categoryUUID)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return mapCategory(category), nil
}

func (r *categoryRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Category, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	categories, err := r.query.GetCategoriesByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Category, len(categories))
	for i, c := range categories {
		result[i] = mapCategory(c)
	}

	return result, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	categoryUUID, err := parseUUID(category.ID)
	if err != nil {
		return domain.ErrNotFound
	}

	updated, err := r.query.UpdateCategory(ctx, db.UpdateCategoryParams{
		ID:    categoryUUID,
		Name:  category.Name,
		Type:  string(category.Type),
		Color: category.Color,
		Icon:  category.Icon,
	})
	if err != nil {
		if isNotFound(err) {
			return domain.ErrNotFound
		}
		return err
	}

	category.UpdatedAt = updated.UpdatedAt.Time
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	categoryUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	return r.query.DeleteCategory(ctx, categoryUUID)
}
