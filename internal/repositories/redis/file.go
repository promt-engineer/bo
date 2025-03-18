package redis

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/pkg/redis"
	"context"
	"fmt"
	"time"

	rd "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	FileCacheKeyPrefix = "file"
)

type fileRepository struct {
	conn *redis.Client
}

func NewFileRepository(conn *redis.Client) *fileRepository {
	return &fileRepository{
		conn: conn,
	}
}

func (r *fileRepository) Get(ctx context.Context, organizationID uuid.UUID, id uuid.UUID) (*entities.File, error) {
	bts, err := r.conn.Get(ctx, fmt.Sprintf("%s:%s:%s", FileCacheKeyPrefix, organizationID.String(), id.String()))
	if err != nil {
		if errors.Is(err, rd.Nil) {
			return nil, e.ErrEntityNotFound
		} else {
			return nil, err
		}
	}

	file := &entities.File{}
	if err = file.Unmarshal(bts); err != nil {
		return nil, err
	}

	return file, nil
}

func (r *fileRepository) Create(ctx context.Context, organizationID uuid.UUID, file *entities.File, expiration time.Duration) error {
	if err := r.conn.Set(ctx, fmt.Sprintf("%s:%s:%s", FileCacheKeyPrefix, organizationID.String(), file.ID.String()), file, expiration); err != nil {
		return err
	}

	return r.conn.HPush(ctx, r.conn.PrepareKey(FileCacheKeyPrefix, organizationID.String()), file.ID.String(), expiration.String())
}

func (r *fileRepository) Update(ctx context.Context, organizationID uuid.UUID, file *entities.File, expiration time.Duration) error {
	return r.conn.Set(ctx, fmt.Sprintf("%s:%s:%s", FileCacheKeyPrefix, organizationID.String(), file.ID.String()), file, expiration)
}

func (r *fileRepository) Find(ctx context.Context, organizationID uuid.UUID) ([]entities.File, error) {
	ids, err := r.conn.HRange(ctx, r.conn.PrepareKey(FileCacheKeyPrefix, organizationID.String()))
	if err != nil {
		if errors.Is(rd.Nil, err) {
			return nil, e.ErrEntityNotFound
		}

		return nil, err
	}

	files := make([]entities.File, 0)

	for _, i := range ids {
		file, err := r.Get(ctx, organizationID, uuid.MustParse(i))
		if err != nil {
			if errors.Is(e.ErrEntityNotFound, err) {
				err = r.conn.HDelete(ctx, r.conn.PrepareKey(FileCacheKeyPrefix, organizationID.String()), i)

				if err != nil {
					return nil, err
				}

				continue
			} else {
				return nil, err
			}
		}

		files = append(files, *file)
	}

	return files, nil
}
