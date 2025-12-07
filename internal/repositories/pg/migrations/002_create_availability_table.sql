-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS availability (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users (id),
  day_of_week INTEGER NOT NULL, -- 0=dom, 1=seg, 2=ter ... 6=sab
  start_time TIME NOT NULL, -- '09:00:00'
  end_time TIME NOT NULL -- '17:00:00'
);

---- create above / drop below ----
 DROP TABLE IF EXISTS availability;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
