-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  slug VARCHAR(255) UNIQUE NOT NULL
);

---- create above / drop below ---- 
DROP TABLE IF EXISTS users;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
