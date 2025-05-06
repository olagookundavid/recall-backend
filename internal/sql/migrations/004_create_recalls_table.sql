-- +goose Up
CREATE TABLE IF NOT EXISTS fda_recalls ( 
    id UUID NOT NULL REFERENCES tracked_product ON DELETE CASCADE, 
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    recall_id TEXT NOT NULL DEFAULT '',
    fda_description TEXT NOT NULL DEFAULT '',
    date_recalled DATE NOT NULL,
    CONSTRAINT unique_id UNIQUE (id, recall_id));

-- +goose Down
DROP TABLE IF EXISTS fda_recalls;
