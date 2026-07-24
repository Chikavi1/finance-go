package domain

import (
	"context"
	"time"
)

type Attachment struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	TransactionID *string   `json:"transaction_id,omitempty"`
	Filename      string    `json:"filename"`
	OriginalName  string    `json:"original_name"`
	MimeType      string    `json:"mime_type"`
	Size          int64     `json:"size"`
	URL           string    `json:"url"`
	CreatedAt     time.Time `json:"created_at"`
}

type AttachmentRepository interface {
	Create(ctx context.Context, att *Attachment) error
	GetByID(ctx context.Context, id string) (*Attachment, error)
	GetByTransactionID(ctx context.Context, transactionID string) ([]*Attachment, error)
	Delete(ctx context.Context, id string) error
}

type StorageService interface {
	Upload(ctx context.Context, filename, contentType string, data []byte) (string, error)
	Delete(ctx context.Context, filename string) error
}
