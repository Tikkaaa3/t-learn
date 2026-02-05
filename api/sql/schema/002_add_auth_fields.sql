-- +goose Up
ALTER TABLE users 
ADD COLUMN api_key TEXT UNIQUE DEFAULT encode(sha256(random()::text::bytea), 'hex'), 
ADD COLUMN role TEXT NOT NULL DEFAULT 'student'; -- 'student' or 'admin'

-- +goose Down
ALTER TABLE users 
DROP COLUMN api_key,
DROP COLUMN role;
