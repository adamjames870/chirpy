-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: RevokeTokensByUser :exec
UPDATE refresh_tokens
SET updated_at = $1, revoked_at = $1
WHERE user_id = $2;

-- name: RevokeTokenByToken :exec
UPDATE refresh_tokens
SET updated_at = $1, revoked_at = $1
WHERE token = $2;

-- name: GetUserFromRefreshToken :one
SELECT user_id, revoked_at
FROM refresh_tokens
WHERE token = $1 ;