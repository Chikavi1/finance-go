-- name: CreateTag :one
INSERT INTO tags (user_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetTagByID :one
SELECT * FROM tags WHERE id = $1;

-- name: GetTagsByUserID :many
SELECT * FROM tags WHERE user_id = $1 ORDER BY name ASC;

-- name: UpdateTag :one
UPDATE tags
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags WHERE id = $1;

-- name: GetTagByName :one
SELECT * FROM tags WHERE user_id = $1 AND name = $2;
