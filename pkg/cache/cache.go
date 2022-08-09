package cache

import "github.com/go-redis/redis/v8"

type Cache interface {
	GetCache() redis.UniversalClient
	Close() error
}
