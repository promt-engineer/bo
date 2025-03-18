package repositories

import (
	"backoffice/internal/entities"
	"context"
	"time"

	"github.com/google/uuid"
)

type FileRepository interface {
	Get(ctx context.Context, organizationID uuid.UUID, id uuid.UUID) (*entities.File, error)
	Create(ctx context.Context, organizationID uuid.UUID, file *entities.File, expiration time.Duration) error
	Update(ctx context.Context, organizationID uuid.UUID, file *entities.File, expiration time.Duration) error
	Find(ctx context.Context, organizationID uuid.UUID) ([]entities.File, error)
}
