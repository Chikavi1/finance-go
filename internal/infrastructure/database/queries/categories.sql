-- name: CreateCategory :one
INSERT INTO categories (user_id, name, type, color, icon)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1;

-- name: GetCategoriesByUserID :many
SELECT * FROM categories WHERE user_id = $1 ORDER BY name ASC;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2,
    type = $3,
    color = $4,
    icon = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;
