-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserBySlug :one
SELECT * FROM users
WHERE slug = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, slug)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;