package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
)

type tagRepository struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func NewTagRepository(pool *pgxpool.Pool) domain.TagRepository {
	return &tagRepository{
		pool:  pool,
		query: db.New(pool),
	}
}

func (r *tagRepository) Create(ctx context.Context, tag *domain.Tag) error {
	userUUID, err := parseUUID(tag.UserID)
	if err != nil {
		return domain.ErrNotFound
	}

	created, err := r.query.CreateTag(ctx, db.CreateTagParams{
		UserID: userUUID,
		Name:   tag.Name,
	})
	if err != nil {
		return err
	}

	tag.ID = pgUUIDToString(created.ID)
	tag.CreatedAt = created.CreatedAt.Time
	return nil
}

func (r *tagRepository) GetByID(ctx context.Context, id string) (*domain.Tag, error) {
	tagUUID, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	tag, err := r.query.GetTagByID(ctx, tagUUID)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return mapTag(tag), nil
}

func (r *tagRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Tag, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	tags, err := r.query.GetTagsByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Tag, len(tags))
	for i, t := range tags {
		result[i] = mapTag(t)
	}

	return result, nil
}

func (r *tagRepository) Update(ctx context.Context, tag *domain.Tag) error {
	tagUUID, err := parseUUID(tag.ID)
	if err != nil {
		return domain.ErrNotFound
	}

	updated, err := r.query.UpdateTag(ctx, db.UpdateTagParams{
		ID:   tagUUID,
		Name: tag.Name,
	})
	if err != nil {
		if isNotFound(err) {
			return domain.ErrNotFound
		}
		return err
	}

	tag.CreatedAt = updated.CreatedAt.Time
	return nil
}

func (r *tagRepository) Delete(ctx context.Context, id string) error {
	tagUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	return r.query.DeleteTag(ctx, tagUUID)
}

func (r *tagRepository) GetOrCreate(ctx context.Context, userID, name string) (*domain.Tag, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	tag, err := r.query.GetTagByName(ctx, db.GetTagByNameParams{
		UserID: userUUID,
		Name:   name,
	})
	if err == nil {
		return mapTag(tag), nil
	}

	if !isNotFound(err) {
		return nil, err
	}

	created, err := r.query.CreateTag(ctx, db.CreateTagParams{
		UserID: userUUID,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}

	return mapTag(created), nil
}
