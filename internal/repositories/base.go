package repositories

import (
	"backoffice/internal/entities"
	"context"
)

type BaseRepository[T any] interface {
	Find(ctx context.Context, conditions map[string]interface{}) (data []*T, err error)
	FindLimit(ctx context.Context, conditions map[string]interface{}, limit, offset int) (data []*T, total int64, err error)
	FindBy(ctx context.Context, params map[string]interface{}) (m *T, err error)
	FindByWith(ctx context.Context, params map[string]interface{}, references ...string) (entity *T, err error)
	Create(ctx context.Context, m *T) (*T, error)
	CreateNoReturn(ctx context.Context, m *T) error
	Delete(ctx context.Context, m *T, conditions ...interface{}) error
	Save(ctx context.Context, m *T) (*T, error)
	Update(ctx context.Context, entity *T, values interface{}, conditions map[string]interface{}) (*T, error)

	Paginate(ctx context.Context, filters map[string]interface{}, order string, limit int, page int) (pagination entities.Pagination[T], err error)
}
