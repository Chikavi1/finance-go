package attachment

import (
	"context"
	"fmt"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	Create(ctx context.Context, userID string, transactionID *string, originalName, mimeType string, data []byte) (*domain.Attachment, error)
	GetByID(ctx context.Context, userID, id string) (*domain.Attachment, error)
	GetByTransactionID(ctx context.Context, userID, transactionID string) ([]*domain.Attachment, error)
	Delete(ctx context.Context, userID, id string) error
}

type service struct {
	attachmentRepo domain.AttachmentRepository
	storageService domain.StorageService
}

func NewService(attachmentRepo domain.AttachmentRepository, storageService domain.StorageService) Service {
	return &service{
		attachmentRepo: attachmentRepo,
		storageService: storageService,
	}
}

func (s *service) Create(ctx context.Context, userID string, transactionID *string, originalName, mimeType string, data []byte) (*domain.Attachment, error) {
	url, err := s.storageService.Upload(ctx, originalName, mimeType, data)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	att := &domain.Attachment{
		UserID:        userID,
		TransactionID: transactionID,
		Filename:      originalName,
		OriginalName:  originalName,
		MimeType:      mimeType,
		Size:          int64(len(data)),
		URL:           url,
	}

	if err := s.attachmentRepo.Create(ctx, att); err != nil {
		return nil, fmt.Errorf("failed to save attachment: %w", err)
	}

	return att, nil
}

func (s *service) GetByID(ctx context.Context, userID, id string) (*domain.Attachment, error) {
	att, err := s.attachmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if att.UserID != userID {
		return nil, domain.ErrNotFound
	}
	return att, nil
}

func (s *service) GetByTransactionID(ctx context.Context, userID, transactionID string) ([]*domain.Attachment, error) {
	attachments, err := s.attachmentRepo.GetByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Attachment, 0, len(attachments))
	for _, a := range attachments {
		if a.UserID == userID {
			result = append(result, a)
		}
	}

	return result, nil
}

func (s *service) Delete(ctx context.Context, userID, id string) error {
	att, err := s.attachmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if att.UserID != userID {
		return domain.ErrNotFound
	}

	if err := s.storageService.Delete(ctx, att.Filename); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return s.attachmentRepo.Delete(ctx, id)
}
