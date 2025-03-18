package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Client struct {
	cfg   *Config
	redis *redis.Client
}

func New(cfg *Config) (c *Client, err error) {
	c = &Client{
		cfg: cfg,
		redis: redis.NewClient(&redis.Options{
			Addr:     cfg.Host + ":" + cfg.Port,
			Password: cfg.Pass,
		}),
	}

	if _, err = c.redis.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.redis.WithContext(ctx).Set(ctx, c.PrepareKey(c.cfg.Prefix, key), value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	return c.redis.WithContext(ctx).Get(ctx, c.PrepareKey(c.cfg.Prefix, key)).Bytes()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	for i, v := range keys {
		keys[i] = c.PrepareKey(c.cfg.Prefix, v)
	}
	return c.redis.WithContext(ctx).Del(ctx, keys...).Err()
}

func (c *Client) PrepareKey(prefix, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

func (c *Client) HPush(ctx context.Context, key string, value interface{}, expiration string) error {
	return c.redis.WithContext(ctx).HSet(ctx, c.PrepareKey(c.cfg.Prefix, key), value, expiration).Err()
}

func (c *Client) HRange(ctx context.Context, key string) ([]string, error) {
	return c.redis.WithContext(ctx).HKeys(ctx, c.PrepareKey(c.cfg.Prefix, key)).Result()
}

func (c *Client) HDelete(ctx context.Context, key, field string) error {
	return c.redis.WithContext(ctx).HDel(ctx, c.PrepareKey(c.cfg.Prefix, key), field).Err()
}
