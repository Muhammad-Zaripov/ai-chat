package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "github.com/udevs/ai-chat/internal/domain/chat"
	sqlcgen "github.com/udevs/ai-chat/internal/infrastructure/postgres/sqlc/gen"
)

type ChatRepository struct {
	q *sqlcgen.Queries
}

func NewChatRepository(pool *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{q: sqlcgen.New(pool)}
}

func (r *ChatRepository) Create(ctx context.Context, c *domain.Chat) (*domain.Chat, error) {
	row, err := r.q.CreateChat(ctx, sqlcgen.CreateChatParams{
		ID:             c.ID,
		Title:          c.Title,
		Model:          c.Model,
		LastResponseID: c.LastResponseID,
		CreatedAt:      c.CreatedAt,
	})
	if err != nil {
		return nil, err
	}
	return chatToDomain(row), nil
}

func (r *ChatRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Chat, error) {
	row, err := r.q.GetChat(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return chatToDomain(row), nil
}

func (r *ChatRepository) List(ctx context.Context, limit, offset int32) ([]*domain.Chat, error) {
	rows, err := r.q.ListChats(ctx, sqlcgen.ListChatsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*domain.Chat, 0, len(rows))
	for _, row := range rows {
		out = append(out, chatToDomain(row))
	}
	return out, nil
}

func (r *ChatRepository) UpdateResponseID(ctx context.Context, id uuid.UUID, responseID string) (*domain.Chat, error) {
	row, err := r.q.UpdateChatResponseID(ctx, sqlcgen.UpdateChatResponseIDParams{
		ID:             id,
		LastResponseID: &responseID,
		UpdatedAt:      time.Now().UTC(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return chatToDomain(row), nil
}

func chatToDomain(row sqlcgen.Chat) *domain.Chat {
	return &domain.Chat{
		ID:             row.ID,
		Title:          row.Title,
		Model:          row.Model,
		LastResponseID: row.LastResponseID,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}
}
