package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
)

type attachmentRepository struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func NewAttachmentRepository(pool *pgxpool.Pool) domain.AttachmentRepository {
	return &attachmentRepository{
		pool:  pool,
		query: db.New(pool),
	}
}

func (r *attachmentRepository) Create(ctx context.Context, att *domain.Attachment) error {
	created, err := r.query.CreateAttachment(ctx, db.CreateAttachmentParams{
		UserID:        mustParseUUID(att.UserID),
		TransactionID: toNullableUUID(att.TransactionID),
		Filename:      att.Filename,
		OriginalName:  att.OriginalName,
		MimeType:      att.MimeType,
		Size:          att.Size,
		Url:           att.URL,
	})
	if err != nil {
		return err
	}

	att.ID = pgUUIDToString(created.ID)
	att.CreatedAt = created.CreatedAt.Time
	return nil
}

func (r *attachmentRepository) GetByID(ctx context.Context, id string) (*domain.Attachment, error) {
	attUUID, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	att, err := r.query.GetAttachmentByID(ctx, attUUID)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return mapAttachment(att), nil
}

func (r *attachmentRepository) GetByTransactionID(ctx context.Context, transactionID string) ([]*domain.Attachment, error) {
	txUUID, err := parseUUID(transactionID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	attachments, err := r.query.GetAttachmentsByTransactionID(ctx, txUUID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Attachment, len(attachments))
	for i, a := range attachments {
		result[i] = mapAttachment(a)
	}

	return result, nil
}

func (r *attachmentRepository) Delete(ctx context.Context, id string) error {
	attUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	return r.query.DeleteAttachment(ctx, attUUID)
}

func mapAttachment(a db.Attachment) *domain.Attachment {
	att := &domain.Attachment{
		ID:           pgUUIDToString(a.ID),
		UserID:       pgUUIDToString(a.UserID),
		Filename:     a.Filename,
		OriginalName: a.OriginalName,
		MimeType:     a.MimeType,
		Size:         a.Size,
		URL:          a.Url,
		CreatedAt:    a.CreatedAt.Time,
	}

	if a.TransactionID.Valid {
		s := pgUUIDToString(a.TransactionID)
		att.TransactionID = &s
	}

	return att
}
