-- +goose Up
CREATE TABLE IF NOT EXISTS fda_recalls ( 
    id UUID NOT NULL REFERENCES tracked_product ON DELETE CASCADE, 
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    date_purchased DATE NOT NULL);

-- +goose Down
DROP TABLE IF EXISTS fda_recalls;
