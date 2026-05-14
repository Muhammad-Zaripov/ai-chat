package chat

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound   = errors.New("chat: not found")
	ErrEmptyInput = errors.New("chat: empty user input")
)

type Chat struct {
	ID             uuid.UUID
	Title          *string
	Model          string
	LastResponseID *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Repository interface {
	Create(ctx context.Context, c *Chat) (*Chat, error)
	Get(ctx context.Context, id uuid.UUID) (*Chat, error)
	List(ctx context.Context, limit, offset int32) ([]*Chat, error)
	UpdateResponseID(ctx context.Context, id uuid.UUID, responseID string) (*Chat, error)
}

// AIClient sends a user turn to the LLM. previousResponseID is the
// `last_response_id` stored on the chat, or "" for the first turn.
// The implementation is responsible for threading conversation state
// (OpenAI Responses API uses `previous_response_id` for this).
type AIClient interface {
	SendMessage(ctx context.Context, model, previousResponseID, userInput string) (AIReply, error)
	GenerateImage(ctx context.Context, req ImageRequest) (ImageReply, error)
}

type AIReply struct {
	ResponseID string
	Output     string
}

type ImageRequest struct {
	Model   string
	Prompt  string
	Size    string
	Quality string
}

type ImageReply struct {
	B64JSON string
	URL     string
}
