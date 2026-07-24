package domain

import (
	"context"
	"time"
)

type Tag struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type TagRepository interface {
	Create(ctx context.Context, tag *Tag) error
	GetByID(ctx context.Context, id string) (*Tag, error)
	GetByUserID(ctx context.Context, userID string) ([]*Tag, error)
	Update(ctx context.Context, tag *Tag) error
	Delete(ctx context.Context, id string) error
	GetOrCreate(ctx context.Context, userID, name string) (*Tag, error)
}
