package message

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	Body      json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(chatID, senderID uuid.UUID, body json.RawMessage) (*Message, error) {
	if chatID == uuid.Nil {
		return nil, ErrInvalidChatID
	}
	if senderID == uuid.Nil {
		return nil, ErrInvalidSenderID
	}
	if len(body) == 0 || !json.Valid(body) {
		return nil, ErrInvalidBody
	}
	now := time.Now().UTC()
	return &Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		SenderID:  senderID,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (m *Message) UpdateBody(body json.RawMessage) error {
	if len(body) == 0 || !json.Valid(body) {
		return ErrInvalidBody
	}
	m.Body = body
	m.UpdatedAt = time.Now().UTC()
	return nil
}
