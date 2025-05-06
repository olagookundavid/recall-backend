-- +goose Up
CREATE TABLE IF NOT EXISTS notification_tokens ( 
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY, 
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    token text NOT NULL,
    date_updated DATE NOT NULL);

-- +goose Down
DROP TABLE IF EXISTS notification_tokens;
