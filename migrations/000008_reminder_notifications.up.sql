ALTER TABLE reminders ADD COLUMN notification_sent_at TIMESTAMPTZ;

CREATE INDEX idx_reminders_notification_sent_at ON reminders(notification_sent_at);
