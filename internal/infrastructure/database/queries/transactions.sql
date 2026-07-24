-- name: CreateTransaction :one
INSERT INTO transactions (user_id, account_id, to_account_id, category_id, type, amount, description, notes, date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetTransactionByID :one
SELECT * FROM transactions WHERE id = $1;

-- name: GetTransactionsByUserID :many
SELECT * FROM transactions
WHERE user_id = $1
ORDER BY date DESC, created_at DESC;

-- name: GetTransactionsByAccountID :many
SELECT * FROM transactions
WHERE account_id = $1 OR to_account_id = $1
ORDER BY date DESC, created_at DESC;

-- name: GetTransactionsByType :many
SELECT * FROM transactions
WHERE user_id = $1 AND type = $2
ORDER BY date DESC, created_at DESC;

-- name: GetTransactionsByDateRange :many
SELECT * FROM transactions
WHERE user_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC, created_at DESC;

-- name: UpdateTransaction :one
UPDATE transactions
SET account_id = $2,
    to_account_id = $3,
    category_id = $4,
    type = $5,
    amount = $6,
    description = $7,
    notes = $8,
    date = $9,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions WHERE id = $1;

-- name: GetUserTransactionsSummary :one
SELECT
    COUNT(*)::bigint as total_count,
    COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0)::double precision as total_income,
    COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)::double precision as total_expense
FROM transactions
WHERE user_id = $1
  AND date >= $2
  AND date <= $3;

-- name: CreateTransactionTag :exec
INSERT INTO transaction_tags (transaction_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: GetTransactionTags :many
SELECT t.id, t.user_id, t.name, t.created_at
FROM tags t
JOIN transaction_tags tt ON t.id = tt.tag_id
WHERE tt.transaction_id = $1;

-- name: DeleteTransactionTags :exec
DELETE FROM transaction_tags WHERE transaction_id = $1;
