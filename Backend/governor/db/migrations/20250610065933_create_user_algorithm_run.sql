-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementEnd

-- +goose StatementBegin
-- Create the user_algorithm_run table
CREATE TABLE user_algorithm_runs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    is_active BOOLEAN DEFAULT TRUE,
    user_algorithm_id UUID NOT NULL,
    start_cron_schedule TEXT,
    end_cron_schedule TEXT,
    order_domain INT NOT NULL DEFAULT 0,
    market TEXT,
    symbol TEXT,
    status INTEGER NOT NULL DEFAULT 0,
    start_time TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    end_time TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    frequency INTEGER,
    portfolio_size INTEGER,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    stopped_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_algorithm_runs_user_algorithm_id
    ON user_algorithm_runs (user_algorithm_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX CONCURRENTLY IF EXISTS idx_user_algorithm_runs_user_algorithm_id;
DROP TABLE user_algorithm_runs;
-- +goose StatementEnd
