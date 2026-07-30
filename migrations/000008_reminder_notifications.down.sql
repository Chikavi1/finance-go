DROP INDEX IF EXISTS idx_reminders_notification_sent_at;
ALTER TABLE reminders DROP COLUMN IF EXISTS notification_sent_at;
