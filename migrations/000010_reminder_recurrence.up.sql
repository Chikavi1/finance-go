ALTER TABLE reminders
ADD COLUMN recurrence_type VARCHAR(20) NOT NULL DEFAULT 'once' CHECK (recurrence_type IN ('once', 'monthly')),
ADD COLUMN day_of_month INT CHECK (day_of_month BETWEEN 1 AND 31);

UPDATE reminders
SET day_of_month = EXTRACT(DAY FROM due_date)::INT
WHERE day_of_month IS NULL;

CREATE INDEX IF NOT EXISTS idx_reminders_recurrence ON reminders(recurrence_type, day_of_month);
