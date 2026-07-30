ALTER TABLE reminders
ADD COLUMN reminder_time TIME NOT NULL DEFAULT '09:00';

CREATE INDEX idx_reminders_due_date_time ON reminders(due_date, reminder_time);
