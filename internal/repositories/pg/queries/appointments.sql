-- name: CreateAppointment :one
INSERT INTO appointments (user_id, client_name, appointment_date, start_time, end_time)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAppointmentsForDate :many
SELECT * FROM appointments
WHERE user_id = $1 AND appointment_date = $2;

-- name of the query that deletes all appointments
-- name: DeleteAllAppointments :exec
DELETE FROM appointments;