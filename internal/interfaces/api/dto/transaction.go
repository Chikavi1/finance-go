package dto

type CreateTransactionRequest struct {
	AccountID   string   `json:"account_id" validate:"required,uuid"`
	ToAccountID *string  `json:"to_account_id,omitempty" validate:"omitempty,uuid"`
	CategoryID  *string  `json:"category_id,omitempty" validate:"omitempty,uuid"`
	Type        string   `json:"type" validate:"required,oneof=income expense transfer informational"`
	Amount      float64  `json:"amount" validate:"gte=0"`
	Description string   `json:"description" validate:"required,max=500"`
	Notes       *string  `json:"notes,omitempty"`
	Date        string   `json:"date" validate:"required"`
	Tags        []string `json:"tags,omitempty"`
}

type UpdateTransactionRequest struct {
	AccountID   string   `json:"account_id" validate:"required,uuid"`
	ToAccountID *string  `json:"to_account_id,omitempty" validate:"omitempty,uuid"`
	CategoryID  *string  `json:"category_id,omitempty" validate:"omitempty,uuid"`
	Type        string   `json:"type" validate:"required,oneof=income expense transfer informational"`
	Amount      float64  `json:"amount" validate:"gte=0"`
	Description string   `json:"description" validate:"required,max=500"`
	Notes       *string  `json:"notes,omitempty"`
	Date        string   `json:"date" validate:"required"`
	Tags        []string `json:"tags,omitempty"`
}

type TransactionResponse struct {
	ID          string   `json:"id"`
	AccountID   string   `json:"account_id"`
	ToAccountID *string  `json:"to_account_id,omitempty"`
	CategoryID  *string  `json:"category_id,omitempty"`
	Type        string   `json:"type"`
	Amount      float64  `json:"amount"`
	Description string   `json:"description"`
	Notes       *string  `json:"notes,omitempty"`
	Date        string   `json:"date"`
	Tags        []string `json:"tags,omitempty"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type TransactionFilterParams struct {
	Type      string `json:"type,omitempty"`
	AccountID string `json:"account_id,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Page      int    `json:"page,omitempty"`
	PerPage   int    `json:"per_page,omitempty"`
}
