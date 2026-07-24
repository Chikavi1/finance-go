package dto

type CreateTagRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type UpdateTagRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type TagResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}
