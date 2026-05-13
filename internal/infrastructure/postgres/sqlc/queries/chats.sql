-- name: CreateChat :one
INSERT INTO chats (id, title, model, last_response_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $5)
RETURNING id, title, model, last_response_id, created_at, updated_at;

-- name: GetChat :one
SELECT id, title, model, last_response_id, created_at, updated_at
FROM chats
WHERE id = $1;

-- name: UpdateChatResponseID :one
UPDATE chats
SET last_response_id = $2,
    updated_at       = $3
WHERE id = $1
RETURNING id, title, model, last_response_id, created_at, updated_at;
