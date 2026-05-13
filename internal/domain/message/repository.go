package message

import (
	"context"

	"github.com/google/uuid"
)

type ListFilter struct {
	ChatID uuid.UUID
	Limit  int32
	Offset int32
}

type Repository interface {
	Create(ctx context.Context, m *Message) (*Message, error)
	Get(ctx context.Context, id uuid.UUID) (*Message, error)
	ListByChat(ctx context.Context, f ListFilter) ([]*Message, error)
	Update(ctx context.Context, m *Message) (*Message, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
