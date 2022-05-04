package lock

import (
	"context"
	"fmt"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisLock struct {
	key      string
	redisCli redis.UniversalClient
	token    string
}

// NewRedisLock 新建redis实现的分布式锁
func NewRedisLock(r redis.UniversalClient, key string) *RedisLock {
	c := &RedisLock{
		key:      genKey(key),
		redisCli: r,
		token:    genToken(),
	}
	return c
}

// Lock 获取锁
func (r *RedisLock) Lock(ctx context.Context, timeout time.Duration) (bool, error) {
	isSet, err := r.redisCli.SetNX(ctx, r.key, r.token, timeout).Result()
	if err == redis.Nil {
		return false, err
	} else if err != nil {
		log.Errorf("acquires the lock err, key: %s, err: %s", r.key, err.Error())
		return false, err
	}
	return isSet, nil
}

// UnLock 释放锁
func (r *RedisLock) UnLock(ctx context.Context) (bool, error) {
	script := "if redis.call('GET',KEYS[1]) == ARGV[1] then return redis.call('DEL',KEYS[1]) else return 0 end"
	ret, err := r.redisCli.Eval(ctx, script, []string{r.key}, r.token).Result()
	if err != nil {
		return false, err
	}
	n, ok := ret.(int64)
	if !ok {
		return false, err
	}
	return n == 1, nil
}

func genKey(key string) string {
	return fmt.Sprintf(RedisLockKey, key)
}
