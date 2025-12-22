package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type Cache[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Set(ctx context.Context, key string, value T, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type RedisCache[T any] struct {
	client *redis.Client
}

func NewRedisCache[T any](client *redis.Client) *RedisCache[T] {
	return &RedisCache[T]{client: client}
}

func (r *RedisCache[T]) Get(ctx context.Context, key string) (T, error) {
	var result T

	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return result, ErrCacheMiss
	}
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(val), &result)

	return result, err
}

func (r *RedisCache[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisCache[T]) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
