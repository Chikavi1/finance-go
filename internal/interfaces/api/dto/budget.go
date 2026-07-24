package dto

type CreateBudgetRequest struct {
	CategoryID string  `json:"category_id" validate:"required,uuid"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	Month      int32   `json:"month" validate:"required,min=1,max=12"`
	Year       int32   `json:"year" validate:"required,min=2020,max=2100"`
}

type UpdateBudgetRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Spent  float64 `json:"spent" validate:"min=0"`
}

type BudgetResponse struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	CategoryID string  `json:"category_id"`
	Amount     float64 `json:"amount"`
	Spent      float64 `json:"spent"`
	Month      int32   `json:"month"`
	Year       int32   `json:"year"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}
