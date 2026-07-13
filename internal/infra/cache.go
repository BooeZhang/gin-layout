package infra

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/samber/lo"

	"gin-layout/config"
)

var ErrAddrEmpty = errors.New("connect redis: addrs is empty")

type RedisClient struct {
	Client redis.UniversalClient
}

func NewRedis(cfg *config.RedisConfig) (*RedisClient, error) {
	addrs := redisClusterAddrs(cfg)
	if len(addrs) == 0 {
		return nil, ErrAddrEmpty
	}

	client, err := newRedisClient(cfg, addrs)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("connect redis ping failed: %w", err)
	}

	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Close() error {
	if r == nil || r.Client == nil {
		return nil
	}
	return r.Client.Close()
}

func (r *RedisClient) Ping(ctx context.Context) error {
	if r == nil || r.Client == nil {
		return errors.New("redis client is nil")
	}
	return r.Client.Ping(ctx).Err()
}

func newRedisClient(cfg *config.RedisConfig, addrs []string) (redis.UniversalClient, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Mode)) {
	case "", "single", "standalone":
		return redis.NewClient(&redis.Options{
			Addr:     addrs[0],
			Password: cfg.Password,
			DB:       cfg.DB,
			PoolSize: cfg.PoolSize,
		}), nil
	case "cluster":
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Password: cfg.Password,
			PoolSize: cfg.PoolSize,
		}), nil
	default:
		return nil, fmt.Errorf("unsupported redis mode: %q", cfg.Mode)
	}
}

func redisClusterAddrs(cfg *config.RedisConfig) []string {
	return lo.FilterMap(cfg.Addrs, func(x string, _ int) (string, bool) {
		addr := strings.TrimSpace(x)
		if addr == "" {
			return "", false
		}
		return addr, true
	})
}
