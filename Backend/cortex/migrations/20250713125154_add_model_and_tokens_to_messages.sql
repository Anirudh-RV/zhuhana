-- Add migration script here
ALTER TABLE messages
ADD COLUMN model TEXT NOT NULL DEFAULT 'unknown';

ALTER TABLE messages
ADD COLUMN tokens INTEGER NOT NULL DEFAULT 0;
