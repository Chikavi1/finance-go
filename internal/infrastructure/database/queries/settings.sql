-- name: GetSettingsByUserID :many
SELECT * FROM settings
WHERE user_id = $1;

-- name: GetSettingByKey :one
SELECT * FROM settings
WHERE user_id = $1 AND key = $2;

-- name: UpsertSetting :one
INSERT INTO settings (user_id, key, value)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, key)
DO UPDATE SET value = $3, updated_at = NOW()
RETURNING *;

-- name: DeleteSetting :exec
DELETE FROM settings
WHERE user_id = $1 AND key = $2;
