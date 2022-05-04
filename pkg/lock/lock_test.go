package lock

import (
	"github.com/BooeZhang/gin-layout/internal/pkg/options"
	"github.com/BooeZhang/gin-layout/store/redis"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	redisC, err := redis.GetRedisFactoryOr(options.NewRedisOptions())
	if err != nil {
		t.Error(err.Error())
	}
	ctx := context.Background()
	l := NewRedisLock(redisC.GetCache(), "test_lock")
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
	redisC, err := redis.GetRedisFactoryOr(options.NewRedisOptions())
	if err != nil {
		t.Error(err.Error())
	}

	t.Run("should lock/unlock success", func(t *testing.T) {
		ctx := context.Background()
		lock1 := NewRedisLock(redisC.GetCache(), "lock2")
		ok, err := lock1.Lock(ctx, 2*time.Second)
		assert.Nil(t, err)
		assert.True(t, ok)

		ok, err = lock1.UnLock(ctx)
		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("should unlock failed", func(t *testing.T) {
		ctx := context.Background()
		lock2 := NewRedisLock(redisC.GetCache(), "lock3")
		ok, err := lock2.Lock(ctx, 2*time.Second)
		assert.Nil(t, err)
		assert.True(t, ok)

		time.Sleep(3 * time.Second)

		ok, err = lock2.UnLock(ctx)
		assert.Nil(t, err)
		assert.False(t, ok)
	})
}
