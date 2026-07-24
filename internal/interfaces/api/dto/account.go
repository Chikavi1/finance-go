package dto

import "github.com/agnathor/finances-go/internal/domain"

type CreateAccountRequest struct {
	Name     string            `json:"name" validate:"required,min=1,max=255"`
	Type     domain.AccountType `json:"type" validate:"required,oneof=cash wallet bank credit_card savings investment"`
	Currency string            `json:"currency" validate:"required,oneof=USD EUR BRL GBP JPY"`
	Balance  float64           `json:"balance"`
	Color    string            `json:"color" validate:"required"`
	Icon     string            `json:"icon" validate:"required"`
}

type UpdateAccountRequest struct {
	Name     string            `json:"name" validate:"required,min=1,max=255"`
	Type     domain.AccountType `json:"type" validate:"required,oneof=cash wallet bank credit_card savings investment"`
	Currency string            `json:"currency" validate:"required,oneof=USD EUR BRL GBP JPY"`
	Color    string            `json:"color" validate:"required"`
	Icon     string            `json:"icon" validate:"required"`
	Archived bool              `json:"archived"`
}

type AccountResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Name      string            `json:"name"`
	Type      domain.AccountType `json:"type"`
	Currency  string            `json:"currency"`
	Balance   float64           `json:"balance"`
	Color     string            `json:"color"`
	Icon      string            `json:"icon"`
	Archived  bool              `json:"archived"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}
