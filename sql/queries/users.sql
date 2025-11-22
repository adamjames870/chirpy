-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUserNameAndEmail :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at = $4
WHERE id = $1
RETURNING *;

-- name: MakeUserChirpyRed :one
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING *;