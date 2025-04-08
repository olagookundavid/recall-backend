-- +goose Up
CREATE TABLE IF NOT EXISTS users ( 
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY, 
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), 
    name text NOT NULL,
    country TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL DEFAULT '',
    url TEXT NOT NULL DEFAULT '',
    date_of_birth DATE NOT NULL DEFAULT '0001-01-01',
    email citext UNIQUE NOT NULL, 
    password_hash bytea NOT NULL,
    isAdmin bool NOT NULL DEFAULT FALSE,
    version integer NOT NULL DEFAULT 1 );
    
-- +goose Down
DROP TABLE IF EXISTS users;