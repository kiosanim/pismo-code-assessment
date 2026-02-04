package lock

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/lock"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisDistributedLockManager struct {
	cacheConnectionData *adapter.CacheConnectionData
	configuration       config.Configuration
	componentName       string
	logger              logger.Logger
}

func NewRedisDistributedLockManager(cacheConnectionData *adapter.CacheConnectionData, log logger.Logger) *RedisDistributedLockManager {
	manager := &RedisDistributedLockManager{
		cacheConnectionData: cacheConnectionData,
		configuration:       config.Configuration{},
		logger:              log,
	}
	manager.componentName = logger.ComponentNameFromStruct(manager)
	return manager
}

// Lock Trying to acquire a lock
func (r *RedisDistributedLockManager) Lock(ctx context.Context, key string) (*lock.Lock, error) {
	lockValue := r.createLockValue()
	ok, err := r.cacheConnectionData.Rdb.SetNX(ctx, key, lockValue, r.configuration.DistributedLock.TTL).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.DistributedLockFailToAcquire
	}
	return &lock.Lock{Key: key, Value: lockValue, Client: r.cacheConnectionData.Rdb}, nil
}

// WaitToLock Waits until waitingTime for acquire a lock, it will retry in intervals of RetryInterval (see config file) until timeout
func (r *RedisDistributedLockManager) WaitToLock(ctx context.Context, key string, waitingTime time.Duration) (*lock.Lock, error) {
	timeout := time.Now().Add(waitingTime)
	for time.Now().Before(timeout) {
		acquiredLock, err := r.Lock(ctx, key)
		if err != nil {
			return nil, err
		} else if acquiredLock == nil {
			return acquiredLock, nil
		}
		select {
		case <-time.After(r.configuration.DistributedLock.RetryInterval):
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return nil, errors.DistributedLockFailToAcquire
}

// Unlock releases a lock using the Lua script suggested by Redis (https://redis.io/docs/latest/commands/set/#patterns)
func (r *RedisDistributedLockManager) Unlock(ctx context.Context, acquiredLock *lock.Lock) error {
	unlockScript := `if redis.call("get",KEYS[1]) == ARGV[1]
					then
						return redis.call("del",KEYS[1])
					else
						return 0
					end`
	script := redis.NewScript(unlockScript)

	_, err := script.Run(ctx, r.cacheConnectionData.Rdb, []string{acquiredLock.Key}, acquiredLock.Value).Result()
	if err != nil {
		return err
	}
	return nil
}

// createLockValue Generate a Lock value based on time.RFC3339Nano
func (r *RedisDistributedLockManager) createLockValue() string {
	return time.Now().Format(time.RFC3339Nano)
}
