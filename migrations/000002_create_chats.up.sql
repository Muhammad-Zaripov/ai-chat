CREATE TABLE IF NOT EXISTS chats (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title             TEXT,
    model             TEXT        NOT NULL,
    last_response_id  TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
