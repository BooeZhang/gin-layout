package redisx

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"gin-layout/config"
)

var (
	r    redis.UniversalClient
	once sync.Once
)

func DialToRedis(cnf *config.Redis) {
	if cnf == nil {
		log.Fatal().Msg("---> [REDIS] configuration files are empty")
	}
	log.Debug().Msg("Creating new Redis connection pool")
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
			PoolSize:     poolSize,
			TLSConfig:    tlsConfig,
		}

		if cnf.MasterName != "" {
			log.Info().Msg("---> [REDIS] Creating sentinel-backed failover client")
			client = redis.NewFailoverClient(redisOption.Failover())
		} else if cnf.EnableCluster {
			log.Info().Msg("---> [REDIS] Creating cluster client")
			client = redis.NewClusterClient(redisOption.Cluster())
		} else {
			log.Info().Msg("---> [REDIS] Creating single-node client")
			client = redis.NewClient(redisOption.Simple())
		}

		pong, err := client.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal().Err(err).Msg("---> [REDIS] redis connect ping failed")
		} else {
			log.Info().Msgf("---> [REDIS] redis connect ping response: %s", pong)
		}
		r = client
	})

	if r == nil {
		log.Fatal().Msgf("---> [REDIS] failed to get redis store: %+v", r)
	}
}

// GetRedis 获取 redis session
func GetRedis() redis.UniversalClient {
	return r
}

func Close() {
	if r != nil {
		err := r.Close()
		if err != nil {
			log.Error().Err(err).Msg("---> [REDIS] redis close failed")
		}
		log.Info().Msg("---> [REDIS] redis closed")
	}
}
