-- +goose Up
CREATE TABLE IF NOT EXISTS tokens ( 
    token text PRIMARY KEY, 
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE, 
    email citext NOT NULL, 
    expiry timestamp(0) with time zone NOT NULL, 
    scope text NOT NULL );

-- +goose Down
DROP TABLE IF EXISTS tokens;
