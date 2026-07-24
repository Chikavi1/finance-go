-- name: CreateGoal :one
INSERT INTO goals (user_id, name, target_amount, current_amount, target_date, icon, color)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetGoalByID :one
SELECT * FROM goals WHERE id = $1;

-- name: GetGoalsByUserID :many
SELECT * FROM goals WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateGoal :one
UPDATE goals
SET name = $2,
    target_amount = $3,
    current_amount = $4,
    target_date = $5,
    icon = $6,
    color = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteGoal :exec
DELETE FROM goals WHERE id = $1;
