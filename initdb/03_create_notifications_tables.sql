\c notification_db;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS notifications(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender UUID NOT NULL,
    receiver UUID NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    read_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS user_shadow (
    id UUID PRIMARY KEY,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);