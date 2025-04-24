-- +goose Up
CREATE TABLE IF NOT EXISTS tracked_product ( 
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY, 
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE, 
    name text NOT NULL, 
    company_name text NOT NULL,
    store_name text NOT NULL,
    country TEXT NOT NULL DEFAULT '',
    category text NOT NULL,
    phone TEXT NOT NULL DEFAULT '',
    url TEXT NOT NULL DEFAULT '',
    date_purchased DATE NOT NULL);

-- +goose Down
DROP TABLE IF EXISTS tracked_product;
