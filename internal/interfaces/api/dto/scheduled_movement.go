package dto

type CreateScheduledMovementRequest struct {
	AccountID   string   `json:"account_id" validate:"required,uuid"`
	CategoryID  *string  `json:"category_id,omitempty" validate:"omitempty,uuid"`
	Type        string   `json:"type" validate:"required,oneof=income expense"`
	Amount      float64  `json:"amount" validate:"required,gt=0"`
	Description string   `json:"description" validate:"required,max=500"`
	Notes       *string  `json:"notes,omitempty"`
	Frequency   string   `json:"frequency" validate:"required,oneof=daily weekdays weekly monthly yearly"`
	StartDate   string   `json:"start_date" validate:"required"`
	NextRunDate *string  `json:"next_run_date,omitempty"`
	EndDate     *string  `json:"end_date,omitempty"`
	Active      *bool    `json:"active,omitempty"`
}

type UpdateScheduledMovementRequest = CreateScheduledMovementRequest

type ScheduledMovementResponse struct {
	ID                string   `json:"id"`
	UserID            string   `json:"user_id"`
	AccountID         string   `json:"account_id"`
	CategoryID        *string  `json:"category_id,omitempty"`
	Type              string   `json:"type"`
	Amount            float64  `json:"amount"`
	Description       string   `json:"description"`
	Notes             *string  `json:"notes,omitempty"`
	Frequency         string   `json:"frequency"`
	StartDate         string   `json:"start_date"`
	NextRunDate       string   `json:"next_run_date"`
	EndDate           *string  `json:"end_date,omitempty"`
	Active            bool     `json:"active"`
	LastGeneratedDate *string  `json:"last_generated_date,omitempty"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}

type GenerateScheduledMovementsResponse struct {
	Created int                   `json:"created"`
	Items   []TransactionResponse `json:"items"`
}
