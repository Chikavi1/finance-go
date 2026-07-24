package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
)

type userRepository struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepository{
		pool:  pool,
		query: db.New(pool),
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	created, err := r.query.CreateUser(ctx, db.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
		AvatarUrl:    toText(user.AvatarURL),
	})
	if err != nil {
		return err
	}

	user.ID = pgUUIDToString(created.ID)
	user.CreatedAt = created.CreatedAt.Time
	user.UpdatedAt = created.UpdatedAt.Time
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	userUUID, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	user, err := r.query.GetUserByID(ctx, userUUID)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return mapUser(user), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.query.GetUserByEmail(ctx, email)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return mapUser(user), nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	userUUID, err := parseUUID(user.ID)
	if err != nil {
		return domain.ErrNotFound
	}

	updated, err := r.query.UpdateUser(ctx, db.UpdateUserParams{
		ID:        userUUID,
		Name:      user.Name,
		AvatarUrl: toText(user.AvatarURL),
	})
	if err != nil {
		return err
	}

	user.UpdatedAt = updated.UpdatedAt.Time
	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	userUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	_, err = r.query.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           userUUID,
		PasswordHash: passwordHash,
	})
	return err
}
