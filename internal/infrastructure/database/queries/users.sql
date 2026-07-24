-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, avatar_url)
VALUES ($1, $2, $3, $4)
RETURNING id, email, password_hash, name, avatar_url, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, email, password_hash, name, avatar_url, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, name, avatar_url, created_at, updated_at
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET name = COALESCE($2, name),
    avatar_url = COALESCE($3, avatar_url),
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, password_hash, name, avatar_url, created_at, updated_at;

-- name: UpdateUserPassword :one
UPDATE users
SET password_hash = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, password_hash, name, avatar_url, created_at, updated_at;
