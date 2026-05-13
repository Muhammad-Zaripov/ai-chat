package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	domain "github.com/udevs/ai-chat/internal/domain/message"
)

type CreateMessageRequest struct {
	ChatID   uuid.UUID       `json:"chat_id" binding:"required"`
	SenderID uuid.UUID       `json:"sender_id" binding:"required"`
	Message  json.RawMessage `json:"message" binding:"required"`
}

type UpdateMessageRequest struct {
	Message json.RawMessage `json:"message" binding:"required"`
}

type MessageResponse struct {
	ID        uuid.UUID       `json:"id"`
	ChatID    uuid.UUID       `json:"chat_id"`
	SenderID  uuid.UUID       `json:"sender_id"`
	Message   json.RawMessage `json:"message"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func FromDomain(m *domain.Message) MessageResponse {
	return MessageResponse{
		ID:        m.ID,
		ChatID:    m.ChatID,
		SenderID:  m.SenderID,
		Message:   m.Body,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromDomainList(items []*domain.Message) []MessageResponse {
	out := make([]MessageResponse, 0, len(items))
	for _, m := range items {
		out = append(out, FromDomain(m))
	}
	return out
}
