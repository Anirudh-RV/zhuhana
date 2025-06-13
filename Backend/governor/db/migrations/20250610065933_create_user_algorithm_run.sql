-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementEnd

-- +goose StatementBegin
-- Create the user_algorithm_run table
CREATE TABLE user_algorithm_run (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    is_active BOOLEAN DEFAULT TRUE,
    user_algorithm_id UUID NOT NULL,
    start_cron_schedule TEXT,
    end_cron_schedule TEXT,
    order_domain INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    stopped_at TIMESTAMPTZ DEFAULT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_algorithm_run;
-- +goose StatementEnd
