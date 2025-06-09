package scheduler

import (
	"context"
	"time"

	"github.com/bsm/redislock"
)

func TryLock(ctx context.Context, key string, ttl time.Duration) (*redislock.Lock, error) {
	return RedisLockObj.Obtain(ctx, key, ttl, nil)
}
