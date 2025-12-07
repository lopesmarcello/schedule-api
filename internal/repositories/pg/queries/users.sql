-- name: CreateUser :one
INSERT INTO users ("name", "email", "password_hash", "slug")
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT id, name, email, password_hash, slug FROM users
WHERE id = $1;

-- name: GetUserBySlug :one
SELECT id, name, email, password_hash, slug FROM users
WHERE slug = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, slug FROM users
WHERE email = $1;
