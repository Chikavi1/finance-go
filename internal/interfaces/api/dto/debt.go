package dto

type CreateDebtRequest struct {
	Name            string  `json:"name" validate:"required,min=1,max=255"`
	TotalAmount     float64 `json:"total_amount" validate:"required,gt=0"`
	RemainingAmount float64 `json:"remaining_amount" validate:"required,gt=0"`
	InterestRate    float64 `json:"interest_rate"`
	DueDate         *string `json:"due_date,omitempty"`
	Status          string  `json:"status" validate:"oneof=active paid overdue"`
	Notes           *string `json:"notes,omitempty"`
}

type UpdateDebtRequest struct {
	Name            string  `json:"name" validate:"required,min=1,max=255"`
	TotalAmount     float64 `json:"total_amount" validate:"required,gt=0"`
	RemainingAmount float64 `json:"remaining_amount" validate:"required,min=0"`
	InterestRate    float64 `json:"interest_rate"`
	DueDate         *string `json:"due_date,omitempty"`
	Status          string  `json:"status" validate:"oneof=active paid overdue"`
	Notes           *string `json:"notes,omitempty"`
}

type DebtResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	Name            string  `json:"name"`
	TotalAmount     float64 `json:"total_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	InterestRate    float64 `json:"interest_rate"`
	DueDate         *string `json:"due_date,omitempty"`
	Status          string  `json:"status"`
	Notes           *string `json:"notes,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type CreateDebtPaymentRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	PaymentDate string  `json:"payment_date" validate:"required"`
	Notes       *string `json:"notes,omitempty"`
}

type DebtPaymentResponse struct {
	ID          string  `json:"id"`
	DebtID      string  `json:"debt_id"`
	Amount      float64 `json:"amount"`
	PaymentDate string  `json:"payment_date"`
	Notes       *string `json:"notes,omitempty"`
	CreatedAt   string  `json:"created_at"`
}
