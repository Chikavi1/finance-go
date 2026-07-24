package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
)

type settingRepository struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func NewSettingRepository(pool *pgxpool.Pool) domain.SettingRepository {
	return &settingRepository{
		pool:  pool,
		query: db.New(pool),
	}
}

func (r *settingRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Setting, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	settings, err := r.query.GetSettingsByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Setting, len(settings))
	for i, s := range settings {
		result[i] = mapSetting(s)
	}

	return result, nil
}

func (r *settingRepository) Upsert(ctx context.Context, setting *domain.Setting) error {
	userUUID, err := parseUUID(setting.UserID)
	if err != nil {
		return domain.ErrNotFound
	}

	created, err := r.query.UpsertSetting(ctx, db.UpsertSettingParams{
		UserID: userUUID,
		Key:    setting.Key,
		Value:  setting.Value,
	})
	if err != nil {
		return err
	}

	setting.ID = pgUUIDToString(created.ID)
	setting.CreatedAt = created.CreatedAt.Time
	setting.UpdatedAt = created.UpdatedAt.Time
	return nil
}

func (r *settingRepository) Delete(ctx context.Context, userID, key string) error {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return domain.ErrNotFound
	}

	return r.query.DeleteSetting(ctx, db.DeleteSettingParams{
		UserID: userUUID,
		Key:    key,
	})
}

func mapSetting(s db.Setting) *domain.Setting {
	return &domain.Setting{
		ID:        pgUUIDToString(s.ID),
		UserID:    pgUUIDToString(s.UserID),
		Key:       s.Key,
		Value:     s.Value,
		CreatedAt: s.CreatedAt.Time,
		UpdatedAt: s.UpdatedAt.Time,
	}
}
