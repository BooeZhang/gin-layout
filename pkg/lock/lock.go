package lock

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// 分布式锁

const (
	RedisLockKey = "lock:%s"
	EtcdLockKey  = "/lock/%s"
)

type Lock interface {
	Lock(ctx context.Context, timeout time.Duration) (bool, error)
	UnLock(ctx context.Context) (bool, error)
}

func genToken() string {
	id, _ := uuid.NewRandom()
	return id.String()
}
