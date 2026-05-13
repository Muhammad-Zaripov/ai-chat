package message

import "errors"

var (
	ErrInvalidChatID   = errors.New("message: invalid chat_id")
	ErrInvalidSenderID = errors.New("message: invalid sender_id")
	ErrInvalidBody     = errors.New("message: body must be non-empty valid json")
	ErrNotFound        = errors.New("message: not found")
)
