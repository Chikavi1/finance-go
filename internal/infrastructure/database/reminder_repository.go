package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/agnathor/finances-go/internal/domain"
)

type reminderRepository struct {
	pool *pgxpool.Pool
}

func NewReminderRepository(pool *pgxpool.Pool) domain.ReminderRepository {
	return &reminderRepository{pool: pool}
}

func (r *reminderRepository) Create(ctx context.Context, reminder *domain.Reminder) error {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO reminders (user_id, title, amount, due_date, reminder_time, recurrence_type, day_of_month, status, related_type, related_id, notes, notification_sent_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`,
		mustParseUUID(reminder.UserID),
		reminder.Title,
		toFloat8(reminder.Amount),
		pgtype.Date{Time: reminder.DueDate, Valid: true},
		toPgTime(reminder.ReminderTime),
		string(reminder.RecurrenceType),
		toNullableInt4(reminder.DayOfMonth),
		string(reminder.Status),
		toText(reminder.RelatedType),
		toNullableUUID(reminder.RelatedID),
		toText(reminder.Notes),
		toNullableTimestamptz(reminder.NotificationSentAt),
	)

	var (
		id        pgtype.UUID
		createdAt pgtype.Timestamptz
		updatedAt pgtype.Timestamptz
	)
	if err := row.Scan(&id, &createdAt, &updatedAt); err != nil {
		return err
	}
	reminder.ID = pgUUIDToString(id)
	reminder.CreatedAt = createdAt.Time
	reminder.UpdatedAt = updatedAt.Time
	return nil
}

func (r *reminderRepository) GetByID(ctx context.Context, id string) (*domain.Reminder, error) {
	reminderUUID, err := parseUUID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, title, amount, due_date, reminder_time, recurrence_type, day_of_month, status, related_type, related_id, notes, notification_sent_at, created_at, updated_at
		FROM reminders
		WHERE id = $1
	`, reminderUUID)

	reminder, err := scanReminder(row)
	if err != nil {
		if isNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return reminder, nil
}

func (r *reminderRepository) GetByUserID(ctx context.Context, userID string, includeDone bool) ([]*domain.Reminder, error) {
	userUUID, err := parseUUID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, title, amount, due_date, reminder_time, recurrence_type, day_of_month, status, related_type, related_id, notes, notification_sent_at, created_at, updated_at
		FROM reminders
		WHERE user_id = $1 AND ($2::boolean OR status = 'pending')
		ORDER BY due_date ASC, created_at DESC
	`, userUUID, includeDone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders := make([]*domain.Reminder, 0)
	for rows.Next() {
		reminder, err := scanReminder(rows)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}
	return reminders, rows.Err()
}

func (r *reminderRepository) GetDueForNotification(ctx context.Context, dueDate time.Time) ([]*domain.Reminder, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, title, amount, due_date, reminder_time, recurrence_type, day_of_month, status, related_type, related_id, notes, notification_sent_at, created_at, updated_at
		FROM reminders
		WHERE status = 'pending'
		  AND (
		    (
		      recurrence_type = 'once'
		      AND notification_sent_at IS NULL
		      AND (
		        due_date < $1
		        OR (due_date = $1 AND reminder_time <= $2)
		      )
		    )
		    OR (
		      recurrence_type = 'monthly'
		      AND due_date <= $1
		      AND (
		        make_date(EXTRACT(YEAR FROM ($1)::date)::INT, EXTRACT(MONTH FROM ($1)::date)::INT, LEAST(day_of_month, EXTRACT(DAY FROM (date_trunc('month', ($1)::date) + INTERVAL '1 month - 1 day'))::INT)) < $1
		        OR (
		          make_date(EXTRACT(YEAR FROM ($1)::date)::INT, EXTRACT(MONTH FROM ($1)::date)::INT, LEAST(day_of_month, EXTRACT(DAY FROM (date_trunc('month', ($1)::date) + INTERVAL '1 month - 1 day'))::INT)) = $1
		          AND reminder_time <= $2
		        )
		      )
		      AND (
		        notification_sent_at IS NULL
		        OR notification_sent_at::date < make_date(EXTRACT(YEAR FROM ($1)::date)::INT, EXTRACT(MONTH FROM ($1)::date)::INT, LEAST(day_of_month, EXTRACT(DAY FROM (date_trunc('month', ($1)::date) + INTERVAL '1 month - 1 day'))::INT))
		      )
		    )
		  )
		ORDER BY due_date ASC, reminder_time ASC, created_at ASC
	`, pgtype.Date{Time: dueDate, Valid: true}, toPgTime(dueDate.Format("15:04")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders := make([]*domain.Reminder, 0)
	for rows.Next() {
		reminder, err := scanReminder(rows)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}
	return reminders, rows.Err()
}

func (r *reminderRepository) MarkNotificationSent(ctx context.Context, id string, sentAt time.Time) error {
	reminderUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	return r.pool.QueryRow(ctx, `
		UPDATE reminders
		SET notification_sent_at = $2,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id
	`, reminderUUID, pgtype.Timestamptz{Time: sentAt, Valid: true}).Scan(&reminderUUID)
}

func (r *reminderRepository) Update(ctx context.Context, reminder *domain.Reminder) error {
	reminderUUID, err := parseUUID(reminder.ID)
	if err != nil {
		return domain.ErrNotFound
	}

	row := r.pool.QueryRow(ctx, `
		UPDATE reminders
		SET title = $2,
		    amount = $3,
		    due_date = $4,
		    reminder_time = $5,
		    recurrence_type = $6,
		    day_of_month = $7,
		    status = $8,
		    related_type = $9,
		    related_id = $10,
		    notes = $11,
		    notification_sent_at = $12,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`,
		reminderUUID,
		reminder.Title,
		toFloat8(reminder.Amount),
		pgtype.Date{Time: reminder.DueDate, Valid: true},
		toPgTime(reminder.ReminderTime),
		string(reminder.RecurrenceType),
		toNullableInt4(reminder.DayOfMonth),
		string(reminder.Status),
		toText(reminder.RelatedType),
		toNullableUUID(reminder.RelatedID),
		toText(reminder.Notes),
		toNullableTimestamptz(reminder.NotificationSentAt),
	)

	var updatedAt pgtype.Timestamptz
	if err := row.Scan(&updatedAt); err != nil {
		if isNotFound(err) {
			return domain.ErrNotFound
		}
		return err
	}
	reminder.UpdatedAt = updatedAt.Time
	return nil
}

func (r *reminderRepository) Delete(ctx context.Context, id string) error {
	reminderUUID, err := parseUUID(id)
	if err != nil {
		return domain.ErrNotFound
	}
	return r.pool.QueryRow(ctx, `DELETE FROM reminders WHERE id = $1 RETURNING id`, reminderUUID).Scan(&reminderUUID)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanReminder(row scanner) (*domain.Reminder, error) {
	var (
		id          pgtype.UUID
		userID      pgtype.UUID
		amount      pgtype.Float8
		dueDate     pgtype.Date
		reminderTime pgtype.Time
		recurrenceType string
		dayOfMonth pgtype.Int4
		relatedType pgtype.Text
		relatedID   pgtype.UUID
		notes       pgtype.Text
		notificationSentAt pgtype.Timestamptz
		createdAt   pgtype.Timestamptz
		updatedAt   pgtype.Timestamptz
		status      string
		reminder    domain.Reminder
	)

	if err := row.Scan(&id, &userID, &reminder.Title, &amount, &dueDate, &reminderTime, &recurrenceType, &dayOfMonth, &status, &relatedType, &relatedID, &notes, &notificationSentAt, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	reminder.ID = pgUUIDToString(id)
	reminder.UserID = pgUUIDToString(userID)
	reminder.Status = domain.ReminderStatus(status)
	if amount.Valid {
		reminder.Amount = &amount.Float64
	}
	reminder.DueDate = dueDate.Time
	reminder.ReminderTime = fromPgTime(reminderTime)
	reminder.RecurrenceType = domain.ReminderRecurrenceType(recurrenceType)
	if dayOfMonth.Valid {
		day := int(dayOfMonth.Int32)
		reminder.DayOfMonth = &day
	}
	reminder.RelatedType = fromText(relatedType)
	if relatedID.Valid {
		s := pgUUIDToString(relatedID)
		reminder.RelatedID = &s
	}
	reminder.Notes = fromText(notes)
	if notificationSentAt.Valid {
		reminder.NotificationSentAt = &notificationSentAt.Time
	}
	reminder.CreatedAt = createdAt.Time
	reminder.UpdatedAt = updatedAt.Time
	return &reminder, nil
}

func toFloat8(n *float64) pgtype.Float8 {
	if n == nil {
		return pgtype.Float8{Valid: false}
	}
	return pgtype.Float8{Float64: *n, Valid: true}
}

func toNullableTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func toPgTime(value string) pgtype.Time {
	parsed, err := time.Parse("15:04", value)
	if err != nil {
		parsed, _ = time.Parse("15:04", "09:00")
	}

	microseconds := int64(parsed.Hour())*60*60*1000000 + int64(parsed.Minute())*60*1000000 + int64(parsed.Second())*1000000
	return pgtype.Time{Microseconds: microseconds, Valid: true}
}

func fromPgTime(value pgtype.Time) string {
	if !value.Valid {
		return "09:00"
	}

	totalSeconds := value.Microseconds / 1000000
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func toNullableInt4(n *int) pgtype.Int4 {
	if n == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*n), Valid: true}
}
