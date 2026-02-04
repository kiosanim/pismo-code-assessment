package repository

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/cache"
	"github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"time"
)

type RedisRepository struct {
	cacheConnectionData *adapter.CacheConnectionData
	componentName       string
	log                 logger.Logger
}

func NewRedisRepository(cacheConnectionData *adapter.CacheConnectionData, log logger.Logger) *RedisRepository {
	repository := &RedisRepository{
		cacheConnectionData: cacheConnectionData,
		log:                 log,
	}
	repository.componentName = logger.ComponentNameFromStruct(repository)
	return repository
}

// Set Store a key/value in cache
func (r *RedisRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	_, err := r.cacheConnectionData.Rdb.Set(ctx, key, value, ttl).Result()
	if err != nil {
		return errors.CacheInsertionError
	}
	return nil
}

// SetNX Save the key/value in cache only if the key not exists
func (r *RedisRepository) SetNX(ctx context.Context, key string, value string, duration time.Duration) bool {
	res, err := r.cacheConnectionData.Rdb.SetNX(ctx, key, value, duration).Result()
	if err != nil {
		return false
	}
	return res
}

// Get Retrieve a key from the cache
func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	res, err := r.cacheConnectionData.Rdb.Get(ctx, key).Result()
	if err != nil {
		return "", errors.CacheNotFoundError
	}
	return res, nil
}

// Del Remove a key
func (r *RedisRepository) Del(ctx context.Context, key string) error {
	_, err := r.cacheConnectionData.Rdb.Del(ctx, key).Result()
	if err != nil {
		return errors.CacheFailedToDeleteError
	}
	return nil
}

// Exists Check if key exists
func (r *RedisRepository) Exists(ctx context.Context, key string) (bool, error) {
	_, err := r.cacheConnectionData.Rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, errors.CacheNotFoundError
	}
	return true, nil
}

// Expire Expires a key from the cache
func (r *RedisRepository) Expire(ctx context.Context, key string, duration time.Duration) error {
	_, err := r.cacheConnectionData.Rdb.Expire(ctx, key, duration).Result()
	if err != nil {
		return errors.CacheFailedToExpireError
	}
	return nil
}

// HSet Save a structure in the cache
func (r *RedisRepository) HSet(ctx context.Context, key string, value any, ttl time.Duration) error {
	_, err := r.cacheConnectionData.Rdb.HSet(ctx, key, value, ttl).Result()
	if err != nil {
		return errors.CacheInsertionError
	}
	return nil
}

// HGet Get a field of a structure from the cache
func (r *RedisRepository) HGet(ctx context.Context, key string, fieldName string) (string, error) {
	res, err := r.cacheConnectionData.Rdb.HGet(ctx, key, fieldName).Result()
	if err != nil {
		return "", errors.CacheNotFoundError
	}
	return res, nil
}

// HGetAll Get a structure from cache
func (r *RedisRepository) HGetAll(ctx context.Context, key string) (any, error) {
	var cachedObject cache.CachedObject
	err := r.cacheConnectionData.Rdb.HGetAll(ctx, key).Scan(&cachedObject)
	if err != nil {
		return nil, errors.CacheNotFoundError
	}
	return cachedObject.Request, nil

}

// HDel Remove a field from a structure stored by key
func (r *RedisRepository) HDel(ctx context.Context, key string, fieldName string) error {
	_, err := r.cacheConnectionData.Rdb.HDel(ctx, key, fieldName).Result()
	if err != nil {
		return errors.CacheFailedToDeleteError
	}
	return nil
}

// HExists Check if a field exists in a structure stored by key
func (r *RedisRepository) HExists(ctx context.Context, key string, fieldName string) (bool, error) {
	_, err := r.cacheConnectionData.Rdb.HExists(ctx, key, fieldName).Result()
	if err != nil {
		return false, errors.CacheNotFoundError
	}
	return true, nil
}

// HExpire Expires a field from a structure identified as key from the cache
func (r *RedisRepository) HExpire(ctx context.Context, key string, fieldName string, duration time.Duration) error {
	_, err := r.cacheConnectionData.Rdb.HExpire(ctx, key, duration, fieldName).Result()
	if err != nil {
		return errors.CacheFailedToExpireError
	}
	return nil
}
