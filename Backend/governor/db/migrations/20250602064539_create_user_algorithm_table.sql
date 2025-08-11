-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementEnd

-- +goose StatementBegin
-- Create the user_algorithm table
CREATE TABLE user_algorithm (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    script_name TEXT NOT NULL,
    script_url TEXT,
    start_cron_schedule TEXT,
    end_cron_schedule TEXT,
    order_domain INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_algorithm_user_id
    ON user_algorithm (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_algorithm_user_id;
DROP TABLE user_algorithm;
-- +goose StatementEnd
