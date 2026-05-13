package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "github.com/udevs/ai-chat/internal/domain/message"
	sqlcgen "github.com/udevs/ai-chat/internal/infrastructure/postgres/sqlc/gen"
)

type MessageRepository struct {
	q *sqlcgen.Queries
}

func NewMessageRepository(pool *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{q: sqlcgen.New(pool)}
}

func (r *MessageRepository) Create(ctx context.Context, m *domain.Message) (*domain.Message, error) {
	row, err := r.q.CreateMessage(ctx, sqlcgen.CreateMessageParams{
		ID:        m.ID,
		ChatID:    m.ChatID,
		SenderID:  m.SenderID,
		Message:   []byte(m.Body),
		CreatedAt: m.CreatedAt,
	})
	if err != nil {
		return nil, err
	}
	return toDomain(row), nil
}

func (r *MessageRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	row, err := r.q.GetMessage(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return toDomain(row), nil
}

func (r *MessageRepository) ListByChat(ctx context.Context, f domain.ListFilter) ([]*domain.Message, error) {
	rows, err := r.q.ListMessagesByChat(ctx, sqlcgen.ListMessagesByChatParams{
		ChatID: f.ChatID,
		Limit:  f.Limit,
		Offset: f.Offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*domain.Message, 0, len(rows))
	for _, row := range rows {
		out = append(out, toDomain(row))
	}
	return out, nil
}

func (r *MessageRepository) Update(ctx context.Context, m *domain.Message) (*domain.Message, error) {
	row, err := r.q.UpdateMessage(ctx, sqlcgen.UpdateMessageParams{
		ID:        m.ID,
		Message:   []byte(m.Body),
		UpdatedAt: m.UpdatedAt,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return toDomain(row), nil
}

func (r *MessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	n, err := r.q.DeleteMessage(ctx, id)
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func toDomain(row sqlcgen.Message) *domain.Message {
	return &domain.Message{
		ID:        row.ID,
		ChatID:    row.ChatID,
		SenderID:  row.SenderID,
		Body:      row.Message,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
