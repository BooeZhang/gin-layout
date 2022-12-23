package lock

import (
	"github.com/BooeZhang/gin-layout/internal/pkg/config"
	"github.com/BooeZhang/gin-layout/pkg/cache"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	cache.ConnectToRedis(config.NewRedisConfig())
	ctx := context.Background()
	l := NewRedisLock(cache.Redis, "test_lock")
	ok, err := l.Lock(ctx, 10*time.Second)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fatal("lock is not ok")
	}

	ok, err = l.UnLock(ctx)
	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Fatal("UnLock is not ok")
	}
}

func TestLockWithTimeout(t *testing.T) {
	cache.ConnectToRedis(config.NewRedisConfig())

	t.Run("should lock/unlock success", func(t *testing.T) {
		ctx := context.Background()
		lock1 := NewRedisLock(cache.Redis, "lock2")
		ok, err := lock1.Lock(ctx, 2*time.Second)
		assert.Nil(t, err)
		assert.True(t, ok)

		ok, err = lock1.UnLock(ctx)
		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("should unlock failed", func(t *testing.T) {
		ctx := context.Background()
		lock2 := NewRedisLock(cache.Redis, "lock3")
		ok, err := lock2.Lock(ctx, 2*time.Second)
		assert.Nil(t, err)
		assert.True(t, ok)

		time.Sleep(3 * time.Second)

		ok, err = lock2.UnLock(ctx)
		assert.Nil(t, err)
		assert.False(t, ok)
	})
}
