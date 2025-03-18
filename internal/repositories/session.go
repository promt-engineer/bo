package repositories

import (
	"backoffice/internal/entities"
	"context"
	"github.com/google/uuid"
	"time"
)

type SessionRepository interface {
	Get(ctx context.Context, key uuid.UUID) (*entities.Session, error)
	GetKeys(ctx context.Context, accountID uuid.UUID) ([]uuid.UUID, error)
	Create(ctx context.Context, key uuid.UUID, session *entities.Session, expiration time.Time) error
	Delete(ctx context.Context, prefix string, key uuid.UUID) error
	Update(ctx context.Context, key uuid.UUID, session *entities.Session) error
	Find(ctx context.Context, id uuid.UUID) ([]*entities.Session, error)
}
