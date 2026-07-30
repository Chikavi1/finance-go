DROP INDEX IF EXISTS idx_reminders_recurrence;
ALTER TABLE reminders
DROP COLUMN IF EXISTS day_of_month,
DROP COLUMN IF EXISTS recurrence_type;
