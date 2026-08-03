package domain

import (
	"context"
	"time"
)

type ReminderStatus string
type ReminderRecurrenceType string

const (
	ReminderStatusPending   ReminderStatus = "pending"
	ReminderStatusDone      ReminderStatus = "done"
	ReminderStatusDismissed ReminderStatus = "dismissed"
)

const (
	ReminderRecurrenceOnce    ReminderRecurrenceType = "once"
	ReminderRecurrenceMonthly ReminderRecurrenceType = "monthly"
)

type Reminder struct {
	ID                 string                 `json:"id"`
	UserID             string                 `json:"user_id"`
	UserEmail          string                 `json:"-"`
	Title              string                 `json:"title"`
	Amount             *float64               `json:"amount,omitempty"`
	DueDate            time.Time              `json:"due_date"`
	ReminderTime       string                 `json:"reminder_time"`
	RecurrenceType     ReminderRecurrenceType `json:"recurrence_type"`
	DayOfMonth         *int                   `json:"day_of_month,omitempty"`
	Status             ReminderStatus         `json:"status"`
	RelatedType        *string                `json:"related_type,omitempty"`
	RelatedID          *string                `json:"related_id,omitempty"`
	Notes              *string                `json:"notes,omitempty"`
	NotificationSentAt *time.Time             `json:"notification_sent_at,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

type ReminderRepository interface {
	Create(ctx context.Context, reminder *Reminder) error
	GetByID(ctx context.Context, id string) (*Reminder, error)
	GetByUserID(ctx context.Context, userID string, includeDone bool) ([]*Reminder, error)
	GetDueForNotification(ctx context.Context, dueDate time.Time) ([]*Reminder, error)
	MarkNotificationSent(ctx context.Context, id string, sentAt time.Time) error
	Update(ctx context.Context, reminder *Reminder) error
	Delete(ctx context.Context, id string) error
}
