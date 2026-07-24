package category

import (
	"context"
	"fmt"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	Create(ctx context.Context, userID string, name string, categoryType domain.CategoryType, color string, icon string) (*domain.Category, error)
	GetByID(ctx context.Context, userID, categoryID string) (*domain.Category, error)
	GetAll(ctx context.Context, userID string) ([]*domain.Category, error)
	Update(ctx context.Context, userID string, categoryID string, name string, categoryType domain.CategoryType, color string, icon string) (*domain.Category, error)
	Delete(ctx context.Context, userID, categoryID string) error
}

type service struct {
	categoryRepo domain.CategoryRepository
}

func NewService(categoryRepo domain.CategoryRepository) Service {
	return &service{categoryRepo: categoryRepo}
}

func (s *service) Create(ctx context.Context, userID string, name string, categoryType domain.CategoryType, color string, icon string) (*domain.Category, error) {
	category := &domain.Category{
		UserID: userID,
		Name:   name,
		Type:   categoryType,
		Color:  color,
		Icon:   icon,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (s *service) GetByID(ctx context.Context, userID, categoryID string) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if category.UserID != userID {
		return nil, domain.ErrNotFound
	}

	return category, nil
}

func (s *service) GetAll(ctx context.Context, userID string) ([]*domain.Category, error) {
	return s.categoryRepo.GetByUserID(ctx, userID)
}

func (s *service) Update(ctx context.Context, userID string, categoryID string, name string, categoryType domain.CategoryType, color string, icon string) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if category.UserID != userID {
		return nil, domain.ErrNotFound
	}

	category.Name = name
	category.Type = categoryType
	category.Color = color
	category.Icon = icon

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

func (s *service) Delete(ctx context.Context, userID, categoryID string) error {
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return err
	}

	if category.UserID != userID {
		return domain.ErrNotFound
	}

	return s.categoryRepo.Delete(ctx, categoryID)
}
