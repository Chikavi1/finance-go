-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, token_hash, expires_at, created_at, revoked;

-- name: GetRefreshToken :one
SELECT id, user_id, token_hash, expires_at, created_at, revoked
FROM refresh_tokens
WHERE id = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked = TRUE
WHERE id = $1;

-- name: RevokeUserRefreshTokens :exec
UPDATE refresh_tokens
SET revoked = TRUE
WHERE user_id = $1 AND revoked = FALSE;
