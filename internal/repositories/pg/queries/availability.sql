-- name: SetAvailability :exec
INSERT INTO availability (
  user_id, day_of_week, start_time, end_time
) VALUES ( $1, $2, $3, $4  );
