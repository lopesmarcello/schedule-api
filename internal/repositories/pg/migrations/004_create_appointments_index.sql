-- Write your migrate up statements here
CREATE INDEX idx_appointments_user_date ON appointments(user_id, appointment_date);

---- create above / drop below ----
DROP INDEX IF EXISTS idx_appointments_user_date;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
