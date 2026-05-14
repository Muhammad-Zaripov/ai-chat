package chat

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	domain "github.com/udevs/ai-chat/internal/domain/chat"
	messagedomain "github.com/udevs/ai-chat/internal/domain/message"
	appimage "github.com/udevs/ai-chat/internal/application/image"
)

type Service struct {
	repo         domain.Repository
	messages     messagedomain.Repository
	ai           domain.AIClient
	imageSvc     *appimage.Service
	defaultModel string
}

var (
	DefaultUserSenderID      = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	DefaultAssistantSenderID = uuid.MustParse("00000000-0000-0000-0000-000000000002")
)

func NewService(repo domain.Repository, messages messagedomain.Repository, ai domain.AIClient, imageSvc *appimage.Service, defaultModel string) *Service {
	return &Service{repo: repo, messages: messages, ai: ai, imageSvc: imageSvc, defaultModel: defaultModel}
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
	
	aiInput := userInput + "\n\n(System instructions: If the user explicitly asks you to generate or draw an image, reply EXACTLY with '$$IMAGE_REQ$$: <english prompt for the image>' and nothing else. If they do not ask for an image, just converse normally.)"
	reply, err := s.ai.SendMessage(ctx, c.Model, prev, aiInput)
	if err != nil {
		return SendOutput{}, err
	}

	replyOutput := reply.Output
	if strings.HasPrefix(replyOutput, "$$IMAGE_REQ$$:") {
		prompt := strings.TrimSpace(strings.TrimPrefix(replyOutput, "$$IMAGE_REQ$$:"))
		out, err := s.imageSvc.Generate(ctx, appimage.GenerateInput{
			Prompt:  prompt,
			Quality: "low",
		})
		if err == nil && out.B64JSON != "" {
			replyOutput = "![Generated Image](data:image/png;base64," + out.B64JSON + ")"
		} else {
			replyOutput = "Kechirasiz, rasmni generatsiya qilishda xatolik yuz berdi."
			if err != nil {
				replyOutput += " " + err.Error()
			}
		}
	}

	updated, err := s.repo.UpdateResponseID(ctx, c.ID, reply.ResponseID)
	if err != nil {
		// AI call succeeded; surface the reply but report persistence error.
		return SendOutput{Chat: c, Reply: replyOutput}, err
	}
	if err := s.storeTurn(ctx, c.ID, DefaultAssistantSenderID, "assistant", replyOutput); err != nil {
		return SendOutput{Chat: updated, Reply: replyOutput}, err
	}
	return SendOutput{Chat: updated, Reply: replyOutput}, nil
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
