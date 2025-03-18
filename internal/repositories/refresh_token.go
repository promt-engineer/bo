package repositories

import (
	"backoffice/internal/entities"
	"context"
	"github.com/google/uuid"
)

type TokenRepository interface {
	Create(ctx context.Context, token *entities.Token) error
	GetByRefresh(ctx context.Context, t string) (token *entities.Token, err error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) (tokens []*entities.Token, err error)
	DeleteByAccess(ctx context.Context, t string) error
}
