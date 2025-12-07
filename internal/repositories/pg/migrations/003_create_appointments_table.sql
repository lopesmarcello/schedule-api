-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS appointments (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id),
  client_name VARCHAR(255) NOT NULL,
  appointment_date DATE NOT NULL,
  start_time TIME NOT NULL,
  end_time TIME NOT NULL,
  status VARCHAR(50) DEFAULT 'scheduled' -- 'scheduled', 'canceled' etc..

);
---- create above / drop below ----
DROP TABLE IF EXISTS appointments;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
