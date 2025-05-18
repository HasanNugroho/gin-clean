package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HasanNugroho/gin-clean/config"
	"github.com/HasanNugroho/gin-clean/pkg/errors"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(config *config.Config) *RedisCache {
	addr := fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Redis.Pass,
		DB:       0,
	})
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	values, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.ErrNotFound.WithMessage("cache key not found")
		}
		return err
	}
	err = json.Unmarshal([]byte(values), value)
	return err
}

func (c *RedisCache) Incr(ctx context.Context, key string) int {
	counts, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0
		}
		return 0
	}
	return int(counts)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, string(bytes), expiration).Err()
}

func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

func (c *RedisCache) Exist(ctx context.Context, key string) (int64, error) {
	return c.client.Exists(ctx, key).Result()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, errors.ErrNotFound.WithMessage("cache key not found")
		}
		return 0, err
	}
	return ttl, nil
}
