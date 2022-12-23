package cache

import (
	"context"
	"crypto/tls"
	"os"
	"sync"
	"time"

	"github.com/BooeZhang/gin-layout/internal/pkg/config"
	"github.com/BooeZhang/gin-layout/pkg/log"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	Redis redis.UniversalClient
	once  sync.Once
)

// ConnectToRedis 连接redis
func ConnectToRedis(cnf *config.RedisConfig) {
	if cnf == nil {
		log.Error("failed to get redis store fatory")
		os.Exit(1)
	}
	log.Debug("Creating new Redis connection pool")
	var (
		tlsConfig *tls.Config
		client    redis.UniversalClient
	)
	once.Do(func() {
		timeout := 5 * time.Second
		if cnf.Timeout > 0 {
			timeout = time.Duration(cnf.Timeout) * time.Second
		}
		// poolSize applies per cluster node and not for the whole cluster.
		poolSize := 500
		if cnf.MaxActive > 0 {
			poolSize = cnf.MaxActive
		}
		if cnf.UseSSL {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: cnf.SSLInsecureSkipVerify,
			}
		}

		redisOption := &redis.UniversalOptions{
			Addrs:        cnf.Addrs,
			MasterName:   cnf.MasterName,
			Password:     cnf.Password,
			DB:           cnf.Database,
			DialTimeout:  timeout,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
			IdleTimeout:  240 * timeout,
			PoolSize:     poolSize,
			TLSConfig:    tlsConfig,
		}

		if cnf.MasterName != "" {
			log.Info("--> [REDIS] Creating sentinel-backed failover client")
			client = redis.NewFailoverClient(redisOption.Failover())
		} else if cnf.EnableCluster {
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
		Redis = client
	})

	if Redis == nil {
		log.Errorf("failed to get redis store fatory, redisFactory: %+v", Redis)
		os.Exit(1)
	}
}
