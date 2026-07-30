DROP INDEX IF EXISTS idx_reminders_due_date_time;
ALTER TABLE reminders DROP COLUMN IF EXISTS reminder_time;
