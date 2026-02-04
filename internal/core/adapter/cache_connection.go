package adapter

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type CacheConnectionData struct {
	Rdb *redis.Client
}

type CacheConnection interface {
	Connect(ctx context.Context) (*CacheConnectionData, error)
}
