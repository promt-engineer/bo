package redis

import (
	"backoffice/internal/entities"
	e "backoffice/internal/errors"
	"backoffice/pkg/redis"
	"context"
	rd "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

const (
	SessionCacheLifetime   = 24 * time.Hour
	SessionCacheKeyPrefix  = "session"
	AccountsCacheKeyPrefix = "accounts"
)

type sessionRepository struct {
	conn *redis.Client
}

func NewSessionRepository(conn *redis.Client) *sessionRepository {
	return &sessionRepository{
		conn: conn,
	}
}

func (r *sessionRepository) Get(ctx context.Context, key uuid.UUID) (*entities.Session, error) {
	bts, err := r.conn.Get(ctx, r.conn.PrepareKey(SessionCacheKeyPrefix, key.String()))
	if err != nil {
		return nil, err
	}

	session := &entities.Session{}
	if err = session.Unmarshal(bts); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *sessionRepository) GetKeys(ctx context.Context, accountID uuid.UUID) ([]uuid.UUID, error) {
	keysStr, err := r.conn.HRange(ctx, r.conn.PrepareKey(AccountsCacheKeyPrefix, accountID.String()))
	if err != nil {
		return nil, err
	}

	keys := []uuid.UUID{}
	for _, sk := range keysStr {
		id, err := uuid.Parse(sk)
		if err != nil {
			return nil, err
		}

		keys = append(keys, id)
	}

	return keys, nil
}

func (r *sessionRepository) Create(ctx context.Context, key uuid.UUID, session *entities.Session, expiration time.Time) error {
	if err := r.conn.Set(ctx, r.conn.PrepareKey(SessionCacheKeyPrefix, key.String()), session, SessionCacheLifetime); err != nil {
		return err
	}

	return r.conn.HPush(ctx, r.conn.PrepareKey(AccountsCacheKeyPrefix, session.Account.ID.String()), key.String(), expiration.String())
}

func (r *sessionRepository) Delete(ctx context.Context, prefix string, key uuid.UUID) error {
	err := r.conn.Del(ctx, r.conn.PrepareKey(prefix, key.String()))

	if errors.Is(rd.Nil, err) {
		return e.ErrEntityNotFound
	}

	return err
}

func (r *sessionRepository) Update(ctx context.Context, key uuid.UUID, session *entities.Session) error {
	return r.conn.Set(ctx, r.conn.PrepareKey(SessionCacheKeyPrefix, key.String()), session, SessionCacheLifetime)
}

func (r *sessionRepository) Find(ctx context.Context, id uuid.UUID) ([]*entities.Session, error) {
	ids, err := r.conn.HRange(ctx, r.conn.PrepareKey(AccountsCacheKeyPrefix, id.String()))
	if err != nil {
		if errors.Is(rd.Nil, err) {
			return nil, e.ErrEntityNotFound
		}

		return nil, err
	}

	sessions := make([]*entities.Session, 0)

	for _, i := range ids {
		session, err := r.Get(ctx, uuid.MustParse(i))
		if err != nil {
			if errors.Is(rd.Nil, err) {
				err = r.conn.HDelete(ctx, r.conn.PrepareKey(AccountsCacheKeyPrefix, id.String()), i)
			}

			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}
