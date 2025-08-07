-- +goose Up
-- +goose StatementBegin
CREATE TABLE notification (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    link TEXT,
    read BOOLEAN NOT NULL DEFAULT FALSE,
    pinned BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_notification_user_id ON notification(user_id);
CREATE INDEX idx_notification_unread ON notification(user_id) WHERE read = FALSE;
CREATE INDEX idx_notification_created_at ON notification(created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notification;
DROP INDEX IF EXISTS idx_notification_user_id;
DROP INDEX IF EXISTS idx_notification_unread;
DROP INDEX IF EXISTS idx_notification_created_at;
-- +goose StatementEnd
