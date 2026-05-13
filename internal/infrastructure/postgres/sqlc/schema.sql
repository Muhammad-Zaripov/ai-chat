CREATE TABLE messages (
    id          UUID PRIMARY KEY,
    chat_id     UUID        NOT NULL,
    sender_id   UUID        NOT NULL,
    message     JSONB       NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL
);
