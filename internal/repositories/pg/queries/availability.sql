-- name: CreateAvailability :one
INSERT INTO availability (user_id, day_of_week, start_time, end_time)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAvailabilityForDay :many
SELECT * FROM availability
WHERE user_id = $1 AND day_of_week = $2;

-- name: DeleteAllAvailabilities :exec
DELETE FROM availability;