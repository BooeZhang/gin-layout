package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/BooeZhang/gin-layout/internal/pkg/options"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/store"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type datastore struct {
	redisCli redis.UniversalClient
}

var (
	redisFactory store.Cache
	once         sync.Once
)

func (ds *datastore) Close() error {
	if ds.redisCli != nil {
		return ds.redisCli.Close()
	}

	return nil
}

func (ds *datastore) GetCache() redis.UniversalClient {
	return ds.redisCli
}

// GetRedisFactoryOr 使用给定的配置创建 redis 工厂。
func GetRedisFactoryOr(opts *options.RedisOptions) (store.Cache, error) {
	if opts == nil && redisFactory == nil {
		return nil, fmt.Errorf("failed to get redis store fatory")
	}
	log.Debug("Creating new Redis connection pool")
	var (
		tlsConfig *tls.Config
		client    redis.UniversalClient
	)
	once.Do(func() {
		timeout := 5 * time.Second
		if opts.Timeout > 0 {
			timeout = time.Duration(opts.Timeout) * time.Second
		}
		// poolSize applies per cluster node and not for the whole cluster.
		poolSize := 500
		if opts.MaxActive > 0 {
			poolSize = opts.MaxActive
		}
		if opts.UseSSL {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: opts.SSLInsecureSkipVerify,
			}
		}

		redisOption := &redis.UniversalOptions{
			Addrs:        opts.Addrs,
			MasterName:   opts.MasterName,
			Password:     opts.Password,
			DB:           opts.Database,
			DialTimeout:  timeout,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
			IdleTimeout:  240 * timeout,
			PoolSize:     poolSize,
			TLSConfig:    tlsConfig,
		}

		if opts.MasterName != "" {
			log.Info("--> [REDIS] Creating sentinel-backed failover client")
			client = redis.NewFailoverClient(redisOption.Failover())
		} else if opts.EnableCluster {
			log.Info("--> [REDIS] Creating cluster client")
			client = redis.NewClusterClient(redisOption.Cluster())
		} else {
			log.Info("--> [REDIS] Creating single-node client")
			client = redis.NewClient(redisOption.Simple())
		}

		pong, err := client.Ping(context.Background()).Result()
		if err != nil {
			log.Error("redis connect ping failed, err:", zap.Any("err", err))
			os.Exit(1)
		} else {
			log.Info("redis connect ping response:", zap.String("pong", pong))
		}
		redisFactory = &datastore{client}
	})

	if redisFactory == nil {
		return nil, fmt.Errorf("failed to get redis store fatory, redisFactory: %+v", redisFactory)
	}

	return redisFactory, nil
}

func GetRedisFactory() store.Cache {
	return redisFactory
}
