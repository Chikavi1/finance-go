package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EnsureReminderSchema(ctx context.Context, pool *pgxpool.Pool) error {
	queries := []string{
		`ALTER TABLE reminders ADD COLUMN IF NOT EXISTS reminder_time TIME NOT NULL DEFAULT '09:00'`,
		`ALTER TABLE reminders ADD COLUMN IF NOT EXISTS recurrence_type VARCHAR(20) NOT NULL DEFAULT 'once'`,
		`ALTER TABLE reminders ADD COLUMN IF NOT EXISTS day_of_month INT`,
		`ALTER TABLE reminders ADD COLUMN IF NOT EXISTS notification_sent_at TIMESTAMPTZ`,
		`CREATE INDEX IF NOT EXISTS idx_reminders_notification_sent_at ON reminders(notification_sent_at)`,
		`CREATE INDEX IF NOT EXISTS idx_reminders_due_date_time ON reminders(due_date, reminder_time)`,
		`CREATE INDEX IF NOT EXISTS idx_reminders_recurrence ON reminders(recurrence_type, day_of_month)`,
		`UPDATE reminders SET day_of_month = EXTRACT(DAY FROM due_date)::INT WHERE day_of_month IS NULL`,
	}

	for _, query := range queries {
		if _, err := pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("ensure reminders schema: %w", err)
		}
	}
	return nil
}
