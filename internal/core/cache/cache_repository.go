package cache

import (
	"context"
	"time"
)

type CachedObject struct {
	Request  any
	Response any
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	SetNX(ctx context.Context, key string, value string, duration time.Duration) bool
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, duration time.Duration) error
	HSet(ctx context.Context, key string, value any, ttl time.Duration) error
	HGet(ctx context.Context, key string, fieldName string) (string, error)
	HGetAll(ctx context.Context, key string) (any, error)
	HDel(ctx context.Context, key string, fieldName string) error
	HExists(ctx context.Context, key string, fieldName string) (bool, error)
	HExpire(ctx context.Context, key string, fieldName string, duration time.Duration) error
}
