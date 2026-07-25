package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
)

type refreshTokenRepository struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) domain.RefreshTokenRepository {
	return &refreshTokenRepository{
		pool:  pool,
		query: db.New(pool),
	}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	tokenUUID, err := parseUUID(token.ID)
	if err != nil {
		return domain.ErrTokenInvalid
	}

	created := db.RefreshToken{}
	err = r.pool.QueryRow(ctx,
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, token_hash, expires_at, created_at, revoked`,
		tokenUUID,
		mustParseUUID(token.UserID),
		token.TokenHash,
		toTimestamptz(token.ExpiresAt),
	).Scan(
		&created.ID,
		&created.UserID,
		&created.TokenHash,
		&created.ExpiresAt,
		&created.CreatedAt,
		&created.Revoked,
	)
	if err != nil {
		return err
	}

	token.ID = pgUUIDToString(created.ID)
	token.CreatedAt = created.CreatedAt.Time
	token.Revoked = created.Revoked
	return nil
}

func (r *refreshTokenRepository) GetByID(ctx context.Context, id string) (*domain.RefreshToken, error) {
	tokenUUID, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	token, err := r.query.GetRefreshToken(ctx, tokenUUID)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return mapRefreshToken(token), nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id string) error {
	tokenUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	return r.query.RevokeRefreshToken(ctx, tokenUUID)
}

func (r *refreshTokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	return r.query.RevokeUserRefreshTokens(ctx, mustParseUUID(userID))
}
