ALTER TABLE reminders
ADD COLUMN IF NOT EXISTS reminder_time TIME NOT NULL DEFAULT '09:00';

CREATE INDEX IF NOT EXISTS idx_reminders_due_date_time ON reminders(due_date, reminder_time);
