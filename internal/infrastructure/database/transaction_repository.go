package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
	db "github.com/agnathor/finances-go/internal/infrastructure/database/db"
)

type transactionRepository struct {
	pool  *pgxpool.Pool
	query *db.Queries
}

func NewTransactionRepository(pool *pgxpool.Pool) domain.TransactionRepository {
	return &transactionRepository{
		pool:  pool,
		query: db.New(pool),
	}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	params := db.CreateTransactionParams{
		UserID:      mustParseUUID(tx.UserID),
		AccountID:   mustParseUUID(tx.AccountID),
		ToAccountID: toNullableUUID(tx.ToAccountID),
		CategoryID:  toNullableUUID(tx.CategoryID),
		Type:        string(tx.Type),
		Amount:      tx.Amount,
		Description: tx.Description,
		Notes:       toText(tx.Notes),
		Date:        pgtype.Date{Time: tx.Date, Valid: true},
	}

	created, err := r.query.CreateTransaction(ctx, params)
	if err != nil {
		return err
	}

	tx.ID = pgUUIDToString(created.ID)
	tx.CreatedAt = created.CreatedAt.Time
	tx.UpdatedAt = created.UpdatedAt.Time

	return r.syncTags(ctx, tx.ID, tx.UserID, tx.Tags)
}

func (r *transactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	txUUID, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	tx, err := r.query.GetTransactionByID(ctx, txUUID)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	tags, err := r.query.GetTransactionTags(ctx, txUUID)
	if err != nil {
		return nil, err
	}

	return mapTransaction(tx, tags), nil
}

func (r *transactionRepository) GetByUserID(ctx context.Context, userID string, filter domain.TransactionFilter) ([]*domain.Transaction, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	var transactions []db.Transaction

	switch {
	case filter.Type != nil:
		transactions, err = r.query.GetTransactionsByType(ctx, db.GetTransactionsByTypeParams{
			UserID: userUUID,
			Type:   string(*filter.Type),
		})
	case filter.StartDate != nil && filter.EndDate != nil:
		transactions, err = r.query.GetTransactionsByDateRange(ctx, db.GetTransactionsByDateRangeParams{
			UserID:    userUUID,
			Date:      pgtype.Date{Time: *filter.StartDate, Valid: true},
			Date_2:    pgtype.Date{Time: *filter.EndDate, Valid: true},
		})
	default:
		transactions, err = r.query.GetTransactionsByUserID(ctx, userUUID)
	}

	if err != nil {
		return nil, err
	}

	result := make([]*domain.Transaction, len(transactions))
	for i, t := range transactions {
		tags, _ := r.query.GetTransactionTags(ctx, t.ID)
		result[i] = mapTransaction(t, tags)
	}

	return result, nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *domain.Transaction) error {
	txUUID, err := parseUUID(tx.ID)
	if err != nil {
		return domain.ErrNotFound
	}

	params := db.UpdateTransactionParams{
		ID:          txUUID,
		AccountID:   mustParseUUID(tx.AccountID),
		ToAccountID: toNullableUUID(tx.ToAccountID),
		CategoryID:  toNullableUUID(tx.CategoryID),
		Type:        string(tx.Type),
		Amount:      tx.Amount,
		Description: tx.Description,
		Notes:       toText(tx.Notes),
		Date:        pgtype.Date{Time: tx.Date, Valid: true},
	}

	updated, err := r.query.UpdateTransaction(ctx, params)
	if err != nil {
		return err
	}

	tx.UpdatedAt = updated.UpdatedAt.Time

	return r.syncTags(ctx, tx.ID, tx.UserID, tx.Tags)
}

func (r *transactionRepository) Delete(ctx context.Context, id string) error {
	txUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	return r.query.DeleteTransaction(ctx, txUUID)
}

func (r *transactionRepository) syncTags(ctx context.Context, transactionID, userID string, tagNames []string) error {
	txUUID := mustParseUUID(transactionID)

	if err := r.query.DeleteTransactionTags(ctx, txUUID); err != nil {
		return err
	}

	for _, name := range tagNames {
		tagUUID, err := r.getOrCreateTag(ctx, userID, name)
		if err != nil {
			return err
		}
		if err := r.query.CreateTransactionTag(ctx, db.CreateTransactionTagParams{
			TransactionID: txUUID,
			TagID:         tagUUID,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (r *transactionRepository) getOrCreateTag(ctx context.Context, userID, name string) (pgtype.UUID, error) {
	userUUID := mustParseUUID(userID)

	tag, err := r.query.GetTagByName(ctx, db.GetTagByNameParams{
		UserID: userUUID,
		Name:   name,
	})
	if err != nil {
		if isNotFound(err) {
			created, err := r.query.CreateTag(ctx, db.CreateTagParams{
				UserID: userUUID,
				Name:   name,
			})
			if err != nil {
				return pgtype.UUID{}, err
			}
			return created.ID, nil
		}
		return pgtype.UUID{}, err
	}

	return tag.ID, nil
}

func mapTransaction(t db.Transaction, tags []db.Tag) *domain.Transaction {
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	tx := &domain.Transaction{
		ID:          pgUUIDToString(t.ID),
		UserID:      pgUUIDToString(t.UserID),
		AccountID:   pgUUIDToString(t.AccountID),
		Type:        domain.TransactionType(t.Type),
		Amount:      t.Amount,
		Description: t.Description,
		Date:        t.Date.Time,
		CreatedAt:   t.CreatedAt.Time,
		UpdatedAt:   t.UpdatedAt.Time,
		Tags:        tagNames,
	}

	if t.ToAccountID.Valid {
		s := pgUUIDToString(t.ToAccountID)
		tx.ToAccountID = &s
	}
	if t.CategoryID.Valid {
		s := pgUUIDToString(t.CategoryID)
		tx.CategoryID = &s
	}
	if t.Notes.Valid {
		tx.Notes = &t.Notes.String
	}

	return tx
}

func toNullableUUID(s *string) pgtype.UUID {
	if s == nil {
		return pgtype.UUID{Valid: false}
	}
	return mustParseUUID(*s)
}
