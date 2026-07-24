package user

import (
	"context"
	"fmt"

	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/pkg/password"
)

type Service interface {
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID, name string, avatarURL *string) (*domain.User, error)
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error
}

type service struct {
	userRepo domain.UserRepository
}

func NewService(userRepo domain.UserRepository) Service {
	return &service{userRepo: userRepo}
}

func (s *service) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) UpdateProfile(ctx context.Context, userID, name string, avatarURL *string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.AvatarURL = avatarURL

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *service) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !password.Verify(currentPassword, user.PasswordHash) {
		return domain.ErrInvalidCredentials
	}

	hashedPass, err := password.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.userRepo.UpdatePassword(ctx, userID, hashedPass)
}
