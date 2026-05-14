package chat

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	domain "github.com/udevs/ai-chat/internal/domain/chat"
	messagedomain "github.com/udevs/ai-chat/internal/domain/message"
)

type Service struct {
	repo         domain.Repository
	messages     messagedomain.Repository
	ai           domain.AIClient
	defaultModel string
}

var (
	DefaultUserSenderID      = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	DefaultAssistantSenderID = uuid.MustParse("00000000-0000-0000-0000-000000000002")
)

func NewService(repo domain.Repository, messages messagedomain.Repository, ai domain.AIClient, defaultModel string) *Service {
	return &Service{repo: repo, messages: messages, ai: ai, defaultModel: defaultModel}
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

func (s *Service) List(ctx context.Context, limit, offset int32) ([]*domain.Chat, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}

type SendInput struct {
	Input    string
	SenderID uuid.UUID
}

type SendOutput struct {
	Chat  *domain.Chat
	Reply string
}

func (s *Service) Send(ctx context.Context, chatID uuid.UUID, in SendInput) (SendOutput, error) {
	userInput := strings.TrimSpace(in.Input)
	if userInput == "" {
		return SendOutput{}, domain.ErrEmptyInput
	}
	if in.SenderID == uuid.Nil {
		in.SenderID = DefaultUserSenderID
	}
	c, err := s.repo.Get(ctx, chatID)
	if err != nil {
		return SendOutput{}, err
	}
	if err := s.storeTurn(ctx, c.ID, in.SenderID, "user", userInput); err != nil {
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
	if err := s.storeTurn(ctx, c.ID, DefaultAssistantSenderID, "assistant", reply.Output); err != nil {
		return SendOutput{Chat: updated, Reply: reply.Output}, err
	}
	return SendOutput{Chat: updated, Reply: reply.Output}, nil
}

func (s *Service) storeTurn(ctx context.Context, chatID, senderID uuid.UUID, role, content string) error {
	body, err := json.Marshal(map[string]string{
		"role":    role,
		"content": content,
	})
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	_, err = s.messages.Create(ctx, &messagedomain.Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		SenderID:  senderID,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	})
	return err
}
