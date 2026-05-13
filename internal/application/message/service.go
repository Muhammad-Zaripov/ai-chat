package message

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	domain "github.com/udevs/ai-chat/internal/domain/message"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
	Body     json.RawMessage
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*domain.Message, error) {
	m, err := domain.New(in.ChatID, in.SenderID, in.Body)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, m)
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	if id == uuid.Nil {
		return nil, domain.ErrNotFound
	}
	return s.repo.Get(ctx, id)
}

type ListInput struct {
	ChatID uuid.UUID
	Limit  int32
	Offset int32
}

func (s *Service) ListByChat(ctx context.Context, in ListInput) ([]*domain.Message, error) {
	if in.ChatID == uuid.Nil {
		return nil, domain.ErrInvalidChatID
	}
	if in.Limit <= 0 || in.Limit > 200 {
		in.Limit = 50
	}
	if in.Offset < 0 {
		in.Offset = 0
	}
	return s.repo.ListByChat(ctx, domain.ListFilter{
		ChatID: in.ChatID,
		Limit:  in.Limit,
		Offset: in.Offset,
	})
}

type UpdateInput struct {
	ID   uuid.UUID
	Body json.RawMessage
}

func (s *Service) Update(ctx context.Context, in UpdateInput) (*domain.Message, error) {
	existing, err := s.repo.Get(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if err := existing.UpdateBody(in.Body); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, existing)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return domain.ErrNotFound
	}
	return s.repo.Delete(ctx, id)
}
