-- name: CreateAccount :one
INSERT INTO accounts (user_id, name, type, currency, balance, color, icon)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM accounts WHERE id = $1;

-- name: GetAccountsByUserID :many
SELECT * FROM accounts WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateAccount :one
UPDATE accounts
SET name = $2,
    type = $3,
    currency = $4,
    color = $5,
    icon = $6,
    archived = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;
