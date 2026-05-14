package image

import (
	"context"
	"strings"

	domain "github.com/udevs/ai-chat/internal/domain/chat"
)

type Service struct {
	ai           domain.AIClient
	defaultModel string
}

func NewService(ai domain.AIClient, defaultModel string) *Service {
	return &Service{ai: ai, defaultModel: defaultModel}
}

type GenerateInput struct {
	Prompt  string
	Model   string
	Size    string
	Quality string
}

type GenerateOutput struct {
	B64JSON string
	URL     string
}

func (s *Service) Generate(ctx context.Context, in GenerateInput) (GenerateOutput, error) {
	prompt := strings.TrimSpace(in.Prompt)
	if prompt == "" {
		return GenerateOutput{}, domain.ErrEmptyInput
	}
	model := strings.TrimSpace(in.Model)
	if model == "" {
		model = s.defaultModel
	}
	size := strings.TrimSpace(in.Size)
	if size == "" {
		size = "1024x1024"
	}
	quality := strings.TrimSpace(in.Quality)
	if quality == "" {
		quality = "auto"
	}

	reply, err := s.ai.GenerateImage(ctx, domain.ImageRequest{
		Model:   model,
		Prompt:  prompt,
		Size:    size,
		Quality: quality,
	})
	if err != nil {
		return GenerateOutput{}, err
	}
	return GenerateOutput{B64JSON: reply.B64JSON, URL: reply.URL}, nil
}
