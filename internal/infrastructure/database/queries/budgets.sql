-- name: CreateBudget :one
INSERT INTO budgets (user_id, category_id, amount, spent, month, year)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetBudgetByID :one
SELECT * FROM budgets WHERE id = $1;

-- name: GetBudgetsByUserID :many
SELECT * FROM budgets WHERE user_id = $1 ORDER BY year DESC, month DESC;

-- name: GetBudgetsByMonthYear :many
SELECT * FROM budgets WHERE user_id = $1 AND month = $2 AND year = $3 ORDER BY created_at DESC;

-- name: UpdateBudget :one
UPDATE budgets
SET amount = $2,
    spent = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteBudget :exec
DELETE FROM budgets WHERE id = $1;
