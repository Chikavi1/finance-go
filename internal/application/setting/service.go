package setting

import (
	"context"
	"fmt"

	"github.com/agnathor/finances-go/internal/domain"
)

type Service interface {
	GetAll(ctx context.Context, userID string) ([]*domain.Setting, error)
	Update(ctx context.Context, userID string, settings map[string]string) ([]*domain.Setting, error)
}

type service struct {
	settingRepo domain.SettingRepository
}

func NewService(settingRepo domain.SettingRepository) Service {
	return &service{settingRepo: settingRepo}
}

func (s *service) GetAll(ctx context.Context, userID string) ([]*domain.Setting, error) {
	return s.settingRepo.GetByUserID(ctx, userID)
}

func (s *service) Update(ctx context.Context, userID string, settings map[string]string) ([]*domain.Setting, error) {
	for key, value := range settings {
		setting := &domain.Setting{
			UserID: userID,
			Key:    key,
			Value:  value,
		}
		if err := s.settingRepo.Upsert(ctx, setting); err != nil {
			return nil, fmt.Errorf("failed to update setting %s: %w", key, err)
		}
	}

	return s.settingRepo.GetByUserID(ctx, userID)
}
