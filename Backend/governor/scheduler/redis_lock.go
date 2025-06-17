package scheduler

import (
	"context"
	"time"

	"github.com/bsm/redislock"
)

func (scs *SchedulerService) TryLock(ctx context.Context, key string, ttl time.Duration) (*redislock.Lock, error) {
	return scs.redisLockObj.Obtain(ctx, key, ttl, nil)
}
