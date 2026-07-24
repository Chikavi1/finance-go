package dto

type CreateGoalRequest struct {
	Name          string   `json:"name" validate:"required,min=1,max=255"`
	TargetAmount  float64  `json:"target_amount" validate:"required,gt=0"`
	CurrentAmount float64  `json:"current_amount" validate:"min=0"`
	TargetDate    *string  `json:"target_date,omitempty"`
	Icon          string   `json:"icon" validate:"required"`
	Color         string   `json:"color" validate:"required"`
}

type UpdateGoalRequest struct {
	Name          string   `json:"name" validate:"required,min=1,max=255"`
	TargetAmount  float64  `json:"target_amount" validate:"required,gt=0"`
	CurrentAmount float64  `json:"current_amount" validate:"min=0"`
	TargetDate    *string  `json:"target_date,omitempty"`
	Icon          string   `json:"icon" validate:"required"`
	Color         string   `json:"color" validate:"required"`
}

type GoalResponse struct {
	ID            string   `json:"id"`
	UserID        string   `json:"user_id"`
	Name          string   `json:"name"`
	TargetAmount  float64  `json:"target_amount"`
	CurrentAmount float64  `json:"current_amount"`
	TargetDate    *string  `json:"target_date,omitempty"`
	Icon          string   `json:"icon"`
	Color         string   `json:"color"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}
