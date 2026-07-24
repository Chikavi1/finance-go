package dto

import "github.com/agnathor/finances-go/internal/domain"

type CreateCategoryRequest struct {
	Name  string            `json:"name" validate:"required,min=1,max=255"`
	Type  domain.CategoryType `json:"type" validate:"required,oneof=income expense"`
	Color string            `json:"color" validate:"required"`
	Icon  string            `json:"icon" validate:"required"`
}

type UpdateCategoryRequest struct {
	Name  string            `json:"name" validate:"required,min=1,max=255"`
	Type  domain.CategoryType `json:"type" validate:"required,oneof=income expense"`
	Color string            `json:"color" validate:"required"`
	Icon  string            `json:"icon" validate:"required"`
}

type CategoryResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Name      string            `json:"name"`
	Type      domain.CategoryType `json:"type"`
	Color     string            `json:"color"`
	Icon      string            `json:"icon"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}
