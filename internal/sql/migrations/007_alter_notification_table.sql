-- +goose Up
ALTER TABLE notification_tokens ADD CONSTRAINT unique_user_id UNIQUE(user_id);

-- +goose Down
ALTER TABLE notification_tokens DROP CONSTRAINT unique_user_id;