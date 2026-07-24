package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/agnathor/finances-go/internal/config"
	"github.com/agnathor/finances-go/internal/domain"
	"github.com/agnathor/finances-go/pkg/jwt"
	"github.com/agnathor/finances-go/pkg/password"
)

type Service interface {
	Register(ctx context.Context, email, pass, name string) (*domain.User, *domain.TokenPair, error)
	Login(ctx context.Context, email, pass string) (*domain.User, *domain.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
	Logout(ctx context.Context, userID, refreshTokenID string) error
}

type service struct {
	userRepo  domain.UserRepository
	refreshRepo domain.RefreshTokenRepository
	jwt       *jwt.Manager
	cfg       config.JWTConfig
}

func NewService(userRepo domain.UserRepository, refreshRepo domain.RefreshTokenRepository, jwt *jwt.Manager, cfg config.JWTConfig) Service {
	return &service{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		jwt:         jwt,
		cfg:         cfg,
	}
}

func (s *service) Register(ctx context.Context, email, pass, name string) (*domain.User, *domain.TokenPair, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, nil, domain.ErrConflict
	}

	hashedPass, err := password.Hash(pass)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: hashedPass,
		Name:         name,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	tokens, err := s.generateTokens(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *service) Login(ctx context.Context, email, pass string) (*domain.User, *domain.TokenPair, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, domain.ErrInvalidCredentials
	}

	if !password.Verify(pass, user.PasswordHash) {
		return nil, nil, domain.ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	userID, tokenID, err := s.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	tokenHash := hashToken(tokenID)

	stored, err := s.refreshRepo.GetByID(ctx, tokenID)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	if stored.TokenHash != tokenHash {
		return nil, domain.ErrTokenInvalid
	}

	if stored.Revoked {
		_ = s.refreshRepo.RevokeAllForUser(ctx, userID)
		return nil, domain.ErrTokenInvalid
	}

	if time.Now().After(stored.ExpiresAt) {
		return nil, domain.ErrTokenExpired
	}

	_ = s.refreshRepo.Revoke(ctx, tokenID)

	tokens, err := s.generateTokens(ctx, userID)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *service) Logout(ctx context.Context, userID, refreshToken string) error {
	_, tokenID, err := s.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return domain.ErrTokenInvalid
	}
	return s.refreshRepo.Revoke(ctx, tokenID)
}

func (s *service) generateTokens(ctx context.Context, userID string) (*domain.TokenPair, error) {
	accessToken, expiresAt, err := s.jwt.GenerateAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, tokenID, _, err := s.jwt.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshExpiry := time.Now().Add(s.cfg.RefreshExpiration)

	rt := &domain.RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(tokenID),
		ExpiresAt: refreshExpiry,
	}

	if err := s.refreshRepo.Create(ctx, rt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func hashToken(tokenID string) string {
	h := sha256.Sum256([]byte(tokenID))
	return hex.EncodeToString(h[:])
}
