package connection

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	configuration       *config.Configuration
	cacheConnectionData *adapter.CacheConnectionData
}

func NewRedisConnection(configuration *config.Configuration) *RedisConnection {
	return &RedisConnection{
		configuration: configuration,
	}
}

func (r *RedisConnection) Connect(ctx context.Context) (*adapter.CacheConnectionData, error) {
	opt, err := redis.ParseURL(r.configuration.Cache.URL)
	if err != nil {
		return nil, errors.CacheConnectionFailedError
	}
	rdb := redis.NewClient(opt)
	err = rdb.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		return nil, errors.CacheConnectionValidationFailedError
	}
	_, err = rdb.Get(ctx, "foo").Result()
	if err != nil {
		return nil, errors.CacheConnectionValidationFailedError
	}
	r.cacheConnectionData = &adapter.CacheConnectionData{
		Rdb: rdb,
	}
	return r.cacheConnectionData, nil
}
