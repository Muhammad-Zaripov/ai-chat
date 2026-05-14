package dto

import (
	"time"

	"github.com/google/uuid"

	domain "github.com/udevs/ai-chat/internal/domain/chat"
)

type CreateChatRequest struct {
	Title *string `json:"title,omitempty"`
	Model string  `json:"model,omitempty" example:"gpt-4o-mini"`
}

type SendMessageRequest struct {
	Input    string     `json:"input" binding:"required" example:"hello, who are you?"`
	SenderID *uuid.UUID `json:"sender_id,omitempty"`
}

type ChatResponse struct {
	ID             uuid.UUID `json:"id"`
	Title          *string   `json:"title,omitempty"`
	Model          string    `json:"model"`
	LastResponseID *string   `json:"last_response_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SendMessageResponse struct {
	Chat  ChatResponse `json:"chat"`
	Reply string       `json:"reply"`
}

type ListChatsResponse struct {
	Items []ChatResponse `json:"items"`
}

func ChatFromDomain(c *domain.Chat) ChatResponse {
	return ChatResponse{
		ID:             c.ID,
		Title:          c.Title,
		Model:          c.Model,
		LastResponseID: c.LastResponseID,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

func ChatsFromDomain(items []*domain.Chat) []ChatResponse {
	out := make([]ChatResponse, 0, len(items))
	for _, c := range items {
		out = append(out, ChatFromDomain(c))
	}
	return out
}
