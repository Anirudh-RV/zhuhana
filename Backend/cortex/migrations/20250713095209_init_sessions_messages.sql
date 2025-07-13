CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP DEFAULT NOW(),
    user_id UUID NOT NULL,
    algorithm_id UUID,
    title TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID REFERENCES sessions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    user_message TEXT NOT NULL,
    system_message TEXT NOT NULL,
    model TEXT NOT NULL DEFAULT 'unknown',
    tokens INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_algorithm ON sessions(user_id, algorithm_id);

CREATE INDEX IF NOT EXISTS idx_messages_session_id ON messages(session_id);
