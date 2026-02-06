package lock

import (
	"context"
	"errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	coreerr "github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/lock"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisDistributedLockManager struct {
	cacheConnectionData *adapter.CacheConnectionData
	configuration       *config.Configuration
	componentName       string
	log                 logger.Logger
}

func NewRedisDistributedLockManager(
	cacheConnectionData *adapter.CacheConnectionData,
	configuration *config.Configuration,
	log logger.Logger) *RedisDistributedLockManager {
	manager := &RedisDistributedLockManager{
		cacheConnectionData: cacheConnectionData,
		configuration:       configuration,
		componentName:       "RedisDistributedLockManager",
		log:                 log,
	}
	manager.componentName = logger.ComponentNameFromStruct(manager)
	return manager
}

// Lock Trying to acquire a lock
func (r *RedisDistributedLockManager) Lock(ctx context.Context, key string, ttl time.Duration) (*lock.Lock, error) {
	lockValue := r.createLockValue()
	ok, err := r.cacheConnectionData.Rdb.SetNX(
		ctx,
		key,
		lockValue,
		ttl,
	).Result()
	if err != nil {
		r.log.Debug(r.componentName+".Lock", err)
		return nil, err
	}
	if !ok {
		r.log.Debug(r.componentName+".Lock", err)
		return nil, coreerr.DistributedLockFailToAcquire
	}
	return &lock.Lock{Key: key, Value: lockValue, Client: r.cacheConnectionData.Rdb}, nil
}

// WaitToLock Waits until waitingTime for acquire a lock, it will retry in intervals of RetryInterval (see config file) until timeout
func (r *RedisDistributedLockManager) WaitToLock(ctx context.Context, key string, ttl time.Duration, waitingTimeMilliseconds time.Duration, retryMilliseconds time.Duration) (*lock.Lock, error) {
	timeout := time.Now().Add(waitingTimeMilliseconds)
	r.log.Debug(r.componentName+".WaitToLock", "status", "Trying to acquire lock...")
	for time.Now().Before(timeout) {
		acquiredLock, err := r.Lock(ctx, key, ttl)
		if err == nil {
			r.log.Debug(r.componentName+".WaitToLock", "Lock acquired:", acquiredLock.Key)
			return acquiredLock, nil
		}
		if !errors.Is(err, coreerr.DistributedLockFailToAcquire) {
			return nil, err
		}
		select {
		case <-time.After(retryMilliseconds):
			continue
		case <-ctx.Done():
			r.log.Debug(r.componentName+".WaitToLock", "err", err)
			return nil, ctx.Err()
		}
	}
	return nil, coreerr.DistributedLockFailToAcquire
}

func (r *RedisDistributedLockManager) WaitToLockUsingDefaultTimeConfiguration(ctx context.Context, key string) (*lock.Lock, error) {
	r.log.Debug(r.componentName + ".WaitToLockUsingDefaultTimeConfiguration")
	waitingTimeMilliseconds := time.Duration(r.configuration.DistributedLock.WaitingTime) * time.Millisecond
	retryMilliseconds := time.Duration(r.configuration.DistributedLock.RetryInterval) * time.Millisecond
	ttl := time.Duration(r.configuration.DistributedLock.TTL) * time.Millisecond
	lck, err := r.WaitToLock(ctx, key, ttl, waitingTimeMilliseconds, retryMilliseconds)
	return lck, err
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
		r.log.Debug(r.componentName+".Unlock", "failed to release lock:", acquiredLock.Key, "err", err)
		return err
	}
	r.log.Debug(r.componentName+".Unlock", "releasing lock:", acquiredLock.Key)
	return nil
}

// createLockValue Generate a Lock value based on time.RFC3339Nano
func (r *RedisDistributedLockManager) createLockValue() string {
	return time.Now().Format(time.RFC3339Nano)
}
