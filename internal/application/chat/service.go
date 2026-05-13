package chat

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	domain "github.com/udevs/ai-chat/internal/domain/chat"
)

type Service struct {
	repo         domain.Repository
	ai           domain.AIClient
	defaultModel string
}

func NewService(repo domain.Repository, ai domain.AIClient, defaultModel string) *Service {
	return &Service{repo: repo, ai: ai, defaultModel: defaultModel}
}

type CreateInput struct {
	Title *string
	Model string // optional, falls back to service default
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*domain.Chat, error) {
	model := strings.TrimSpace(in.Model)
	if model == "" {
		model = s.defaultModel
	}
	now := time.Now().UTC()
	c := &domain.Chat{
		ID:        uuid.New(),
		Title:     in.Title,
		Model:     model,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.repo.Create(ctx, c)
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	if id == uuid.Nil {
		return nil, domain.ErrNotFound
	}
	return s.repo.Get(ctx, id)
}

type SendOutput struct {
	Chat  *domain.Chat
	Reply string
}

func (s *Service) Send(ctx context.Context, chatID uuid.UUID, userInput string) (SendOutput, error) {
	if strings.TrimSpace(userInput) == "" {
		return SendOutput{}, domain.ErrEmptyInput
	}
	c, err := s.repo.Get(ctx, chatID)
	if err != nil {
		return SendOutput{}, err
	}
	prev := ""
	if c.LastResponseID != nil {
		prev = *c.LastResponseID
	}
	reply, err := s.ai.SendMessage(ctx, c.Model, prev, userInput)
	if err != nil {
		return SendOutput{}, err
	}
	updated, err := s.repo.UpdateResponseID(ctx, c.ID, reply.ResponseID)
	if err != nil {
		// AI call succeeded; surface the reply but report persistence error.
		return SendOutput{Chat: c, Reply: reply.Output}, err
	}
	return SendOutput{Chat: updated, Reply: reply.Output}, nil
}
