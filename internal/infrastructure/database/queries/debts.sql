-- name: CreateDebt :one
INSERT INTO debts (user_id, name, total_amount, remaining_amount, interest_rate, due_date, status, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetDebtByID :one
SELECT * FROM debts WHERE id = $1;

-- name: GetDebtsByUserID :many
SELECT * FROM debts WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateDebt :one
UPDATE debts
SET name = $2,
    total_amount = $3,
    remaining_amount = $4,
    interest_rate = $5,
    due_date = $6,
    status = $7,
    notes = $8,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDebt :exec
DELETE FROM debts WHERE id = $1;

-- name: CreateDebtPayment :one
INSERT INTO debt_payments (debt_id, amount, payment_date, notes)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDebtPaymentsByDebtID :many
SELECT * FROM debt_payments WHERE debt_id = $1 ORDER BY payment_date DESC;

-- name: DeleteDebtPayment :exec
DELETE FROM debt_payments WHERE id = $1;
