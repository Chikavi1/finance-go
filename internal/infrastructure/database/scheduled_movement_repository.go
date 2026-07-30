package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
)

type scheduledMovementRepository struct {
	pool *pgxpool.Pool
}

func NewScheduledMovementRepository(pool *pgxpool.Pool) domain.ScheduledMovementRepository {
	return &scheduledMovementRepository{pool: pool}
}

func (r *scheduledMovementRepository) Create(ctx context.Context, movement *domain.ScheduledMovement) error {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO scheduled_movements (
			user_id, account_id, category_id, type, amount, description, notes,
			frequency, start_date, next_run_date, end_date, active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`,
		mustParseUUID(movement.UserID),
		mustParseUUID(movement.AccountID),
		toNullableUUID(movement.CategoryID),
		string(movement.Type),
		movement.Amount,
		movement.Description,
		toText(movement.Notes),
		string(movement.Frequency),
		pgtype.Date{Time: movement.StartDate, Valid: true},
		pgtype.Date{Time: movement.NextRunDate, Valid: true},
		toNullableDate(movement.EndDate),
		movement.Active,
	)

	var (
		id        pgtype.UUID
		createdAt pgtype.Timestamptz
		updatedAt pgtype.Timestamptz
	)
	if err := row.Scan(&id, &createdAt, &updatedAt); err != nil {
		return err
	}
	movement.ID = pgUUIDToString(id)
	movement.CreatedAt = createdAt.Time
	movement.UpdatedAt = updatedAt.Time
	return nil
}

func (r *scheduledMovementRepository) GetByID(ctx context.Context, id string) (*domain.ScheduledMovement, error) {
	movementUUID, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	row := r.pool.QueryRow(ctx, selectScheduledMovementSQL()+` WHERE id = $1`, movementUUID)
	movement, err := scanScheduledMovement(row)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return movement, nil
}

func (r *scheduledMovementRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.ScheduledMovement, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	rows, err := r.pool.Query(ctx, selectScheduledMovementSQL()+` WHERE user_id = $1 ORDER BY active DESC, next_run_date ASC`, userUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanScheduledMovementRows(rows)
}

func (r *scheduledMovementRepository) GetDueByUserID(ctx context.Context, userID string, dueDate time.Time) ([]*domain.ScheduledMovement, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	rows, err := r.pool.Query(ctx, selectScheduledMovementSQL()+`
		WHERE user_id = $1
		  AND active = TRUE
		  AND next_run_date <= $2
		  AND (end_date IS NULL OR next_run_date <= end_date)
		ORDER BY next_run_date ASC
	`, userUUID, pgtype.Date{Time: dueDate, Valid: true})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanScheduledMovementRows(rows)
}

func (r *scheduledMovementRepository) Update(ctx context.Context, movement *domain.ScheduledMovement) error {
	movementUUID, err := parseUUID(movement.ID)
	if err != nil {
		return domain.ErrNotFound
	}

	row := r.pool.QueryRow(ctx, `
		UPDATE scheduled_movements
		SET account_id = $2,
		    category_id = $3,
		    type = $4,
		    amount = $5,
		    description = $6,
		    notes = $7,
		    frequency = $8,
		    start_date = $9,
		    next_run_date = $10,
		    end_date = $11,
		    active = $12,
		    last_generated_date = $13,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`,
		movementUUID,
		mustParseUUID(movement.AccountID),
		toNullableUUID(movement.CategoryID),
		string(movement.Type),
		movement.Amount,
		movement.Description,
		toText(movement.Notes),
		string(movement.Frequency),
		pgtype.Date{Time: movement.StartDate, Valid: true},
		pgtype.Date{Time: movement.NextRunDate, Valid: true},
		toNullableDate(movement.EndDate),
		movement.Active,
		toNullableDate(movement.LastGeneratedDate),
	)

	var updatedAt pgtype.Timestamptz
	if err := row.Scan(&updatedAt); err != nil {
		if isNotFound(err) {
			return domain.ErrNotFound
		}
		return err
	}
	movement.UpdatedAt = updatedAt.Time
	return nil
}

func (r *scheduledMovementRepository) Delete(ctx context.Context, id string) error {
	movementUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}
	return r.pool.QueryRow(ctx, `DELETE FROM scheduled_movements WHERE id = $1 RETURNING id`, movementUUID).Scan(&movementUUID)
}

func selectScheduledMovementSQL() string {
	return `
		SELECT id, user_id, account_id, category_id, type, amount, description, notes,
		       frequency, start_date, next_run_date, end_date, active, last_generated_date,
		       created_at, updated_at
		FROM scheduled_movements`
}

func scanScheduledMovementRows(rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}) ([]*domain.ScheduledMovement, error) {
	movements := make([]*domain.ScheduledMovement, 0)
	for rows.Next() {
		movement, err := scanScheduledMovement(rows)
		if err != nil {
			return nil, err
		}
		movements = append(movements, movement)
	}
	return movements, rows.Err()
}

func scanScheduledMovement(row scanner) (*domain.ScheduledMovement, error) {
	var (
		id                pgtype.UUID
		userID            pgtype.UUID
		accountID         pgtype.UUID
		categoryID        pgtype.UUID
		notes             pgtype.Text
		startDate         pgtype.Date
		nextRunDate       pgtype.Date
		endDate           pgtype.Date
		lastGeneratedDate pgtype.Date
		createdAt         pgtype.Timestamptz
		updatedAt         pgtype.Timestamptz
		movementType      string
		frequency         string
		movement          domain.ScheduledMovement
	)

	if err := row.Scan(
		&id,
		&userID,
		&accountID,
		&categoryID,
		&movementType,
		&movement.Amount,
		&movement.Description,
		&notes,
		&frequency,
		&startDate,
		&nextRunDate,
		&endDate,
		&movement.Active,
		&lastGeneratedDate,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}

	movement.ID = pgUUIDToString(id)
	movement.UserID = pgUUIDToString(userID)
	movement.AccountID = pgUUIDToString(accountID)
	movement.Type = domain.TransactionType(movementType)
	movement.Frequency = domain.ScheduledMovementFrequency(frequency)
	if categoryID.Valid {
		s := pgUUIDToString(categoryID)
		movement.CategoryID = &s
	}
	movement.Notes = fromText(notes)
	movement.StartDate = startDate.Time
	movement.NextRunDate = nextRunDate.Time
	if endDate.Valid {
		movement.EndDate = &endDate.Time
	}
	if lastGeneratedDate.Valid {
		movement.LastGeneratedDate = &lastGeneratedDate.Time
	}
	movement.CreatedAt = createdAt.Time
	movement.UpdatedAt = updatedAt.Time
	return &movement, nil
}
