-- +goose Up
CREATE TABLE pot_fda_recalls (
    id UUID REFERENCES tracked_product ON DELETE CASCADE, 
    status TEXT NOT NULL DEFAULT '',
    country TEXT NOT NULL DEFAULT '',
    product_type TEXT NOT NULL DEFAULT '',
    recalling_firm TEXT NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    voluntary_mandated TEXT NOT NULL DEFAULT '',
    initial_firm_notification TEXT NOT NULL DEFAULT '',
    distribution_pattern TEXT NOT NULL DEFAULT '',
    recall_number TEXT NOT NULL,
    product_description TEXT NOT NULL DEFAULT '',
    product_quantity TEXT NOT NULL DEFAULT '',
    reason_for_recall TEXT NOT NULL DEFAULT '',
    recall_initiation_date TEXT NOT NULL DEFAULT '',
    termination_date TEXT NOT NULL DEFAULT '',
    report_date TEXT NOT NULL DEFAULT '',
    code_info TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_recall_number UNIQUE (id, recall_number)
);

-- +goose Down
DROP TABLE IF EXISTS pot_fda_recalls;