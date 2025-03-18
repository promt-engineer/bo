package services

import (
	"backoffice/internal/entities"
	"backoffice/pkg/redis"
	"context"
	"errors"
	"fmt"
	r "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"
)

const (
	AccountsCacheLifetime  = 720 * time.Hour
	AccountsCacheKeyPrefix = "accounts"
)

type CacheService struct {
	accountService *AccountService
	cache          *redis.Client
}

func NewCacheService(cache *redis.Client, accountService *AccountService) *CacheService {
	return &CacheService{
		cache:          cache,
		accountService: accountService,
	}
}

func (s *CacheService) GetAccount(ctx context.Context, id uuid.UUID) (*entities.Account, error) {
	account := &entities.Account{}
	bts, err := s.cache.Get(ctx, s.prepareKey(AccountsCacheKeyPrefix, id.String()))
	if err == nil {
		if err = account.Unmarshal(bts); err != nil {
			return nil, err
		}

		return account, nil
	}

	if !errors.Is(err, r.Nil) {
		return nil, err
	}

	account, err = s.accountService.FindBy(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}

	if err = s.cache.Set(ctx, s.prepareKey(AccountsCacheKeyPrefix, id.String()), account, AccountsCacheLifetime); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *CacheService) StoreAccount(ctx context.Context, account *entities.Account) error {
	return s.cache.Set(ctx, s.prepareKey(AccountsCacheKeyPrefix, account.ID.String()), account, AccountsCacheLifetime)
}

func (s *CacheService) DeleteAccount(ctx context.Context, account *entities.Account) error {
	return s.cache.Del(ctx, s.prepareKey(AccountsCacheKeyPrefix, account.ID.String()))
}

func (s *CacheService) prepareKey(prefix, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}
