package dto

type CreateReminderRequest struct {
	Title       string   `json:"title" validate:"required,max=255"`
	Amount      *float64 `json:"amount,omitempty" validate:"omitempty,gt=0"`
	DueDate     string   `json:"due_date" validate:"required"`
	ReminderTime string  `json:"reminder_time,omitempty"`
	RecurrenceType string `json:"recurrence_type,omitempty" validate:"omitempty,oneof=once monthly"`
	DayOfMonth *int `json:"day_of_month,omitempty" validate:"omitempty,gte=1,lte=31"`
	Status      string   `json:"status,omitempty" validate:"omitempty,oneof=pending done dismissed"`
	RelatedType *string  `json:"related_type,omitempty"`
	RelatedID   *string  `json:"related_id,omitempty" validate:"omitempty,uuid"`
	Notes       *string  `json:"notes,omitempty"`
}

type UpdateReminderRequest = CreateReminderRequest

type ReminderResponse struct {
	ID          string   `json:"id"`
	UserID      string   `json:"user_id"`
	Title       string   `json:"title"`
	Amount      *float64 `json:"amount,omitempty"`
	DueDate     string   `json:"due_date"`
	ReminderTime string  `json:"reminder_time"`
	RecurrenceType string `json:"recurrence_type"`
	DayOfMonth *int `json:"day_of_month,omitempty"`
	Status      string   `json:"status"`
	RelatedType *string  `json:"related_type,omitempty"`
	RelatedID   *string  `json:"related_id,omitempty"`
	Notes       *string  `json:"notes,omitempty"`
	NotificationSentAt *string `json:"notification_sent_at,omitempty"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}
