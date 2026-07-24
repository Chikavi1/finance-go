package tag

import (
	"context"
	"fmt"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	Create(ctx context.Context, userID string, name string) (*domain.Tag, error)
	GetByID(ctx context.Context, userID, tagID string) (*domain.Tag, error)
	GetAll(ctx context.Context, userID string) ([]*domain.Tag, error)
	Update(ctx context.Context, userID string, tagID string, name string) (*domain.Tag, error)
	Delete(ctx context.Context, userID, tagID string) error
	GetOrCreate(ctx context.Context, userID, name string) (*domain.Tag, error)
}

type service struct {
	tagRepo domain.TagRepository
}

func NewService(tagRepo domain.TagRepository) Service {
	return &service{tagRepo: tagRepo}
}

func (s *service) Create(ctx context.Context, userID string, name string) (*domain.Tag, error) {
	tag := &domain.Tag{
		UserID: userID,
		Name:   name,
	}

	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}

func (s *service) GetByID(ctx context.Context, userID, tagID string) (*domain.Tag, error) {
	tag, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, err
	}

	if tag.UserID != userID {
		return nil, domain.ErrNotFound
	}

	return tag, nil
}

func (s *service) GetAll(ctx context.Context, userID string) ([]*domain.Tag, error) {
	return s.tagRepo.GetByUserID(ctx, userID)
}

func (s *service) Update(ctx context.Context, userID string, tagID string, name string) (*domain.Tag, error) {
	tag, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, err
	}

	if tag.UserID != userID {
		return nil, domain.ErrNotFound
	}

	tag.Name = name

	if err := s.tagRepo.Update(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return tag, nil
}

func (s *service) Delete(ctx context.Context, userID, tagID string) error {
	tag, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return err
	}

	if tag.UserID != userID {
		return domain.ErrNotFound
	}

	return s.tagRepo.Delete(ctx, tagID)
}

func (s *service) GetOrCreate(ctx context.Context, userID, name string) (*domain.Tag, error) {
	return s.tagRepo.GetOrCreate(ctx, userID, name)
}
