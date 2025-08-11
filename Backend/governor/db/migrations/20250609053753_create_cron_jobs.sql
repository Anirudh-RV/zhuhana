-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementEnd

-- +goose StatementBegin
-- Create the user_algorithm table
CREATE TABLE cron_job (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_algorithm_id UUID,
    cron_entry_id Text,
    schedule TEXT NOT NULL,         -- the cron expression (e.g. "* * * * 1-3")
    job_type TEXT NOT NULL,
    kafka_topic TEXT NOT NULL,      -- or any payload metadata
    is_active BOOLEAN DEFAULT TRUE, -- controls whether job is live
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cron_job_user_algorithm_id
    ON cron_job (user_algorithm_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_cron_job_user_algorithm_id;
DROP TABLE cron_job;
-- +goose StatementEnd
