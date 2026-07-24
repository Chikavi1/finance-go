-- name: CreateAttachment :one
INSERT INTO attachments (user_id, transaction_id, filename, original_name, mime_type, size, url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAttachmentByID :one
SELECT * FROM attachments WHERE id = $1;

-- name: GetAttachmentsByTransactionID :many
SELECT * FROM attachments WHERE transaction_id = $1 ORDER BY created_at DESC;

-- name: DeleteAttachment :exec
DELETE FROM attachments WHERE id = $1;
