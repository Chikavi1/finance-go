package domain

import (
	"context"
	"time"
)

type Setting struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SettingRepository interface {
	GetByUserID(ctx context.Context, userID string) ([]*Setting, error)
	Upsert(ctx context.Context, setting *Setting) error
	Delete(ctx context.Context, userID, key string) error
}
