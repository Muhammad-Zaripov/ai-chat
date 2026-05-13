-- name: CreateMessage :one
INSERT INTO messages (id, chat_id, sender_id, message, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $5)
RETURNING id, chat_id, sender_id, message, created_at, updated_at;

-- name: GetMessage :one
SELECT id, chat_id, sender_id, message, created_at, updated_at
FROM messages
WHERE id = $1;

-- name: ListMessagesByChat :many
SELECT id, chat_id, sender_id, message, created_at, updated_at
FROM messages
WHERE chat_id = $1
ORDER BY created_at ASC, id ASC
LIMIT $2 OFFSET $3;

-- name: UpdateMessage :one
UPDATE messages
SET message    = $2,
    updated_at = $3
WHERE id = $1
RETURNING id, chat_id, sender_id, message, created_at, updated_at;

-- name: DeleteMessage :execrows
DELETE FROM messages
WHERE id = $1;
