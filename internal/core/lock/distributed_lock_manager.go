package lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	AccountCreationLockKey     = "lock-account-creation"
	TransactionCreationLockKey = "lock-transaction-creation"
)

type Lock struct {
	Key    string        `redis:"key"`
	Value  string        `redis:"value"`
	Client *redis.Client `redis:"-"`
}

type DistributedLockManager interface {
	Lock(ctx context.Context, key string) (*Lock, error)
	WaitToLock(ctx context.Context, key string, waitingTime time.Duration) (*Lock, error)
	Unlock(ctx context.Context, acquiredLock *Lock) error
}
