CREATE TABLE messages (
    id          UUID PRIMARY KEY,
    chat_id     UUID        NOT NULL,
    sender_id   UUID        NOT NULL,
    message     JSONB       NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL
);

CREATE TABLE chats (
    id                UUID PRIMARY KEY,
    title             TEXT,
    model             TEXT        NOT NULL,
    last_response_id  TEXT,
    created_at        TIMESTAMPTZ NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL
);
