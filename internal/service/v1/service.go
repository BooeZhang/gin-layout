package v1

import (
	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/store"
	"github.com/go-redis/redis/v8"
)

type ServiceContext struct {
	Rs  redis.UniversalClient
	Cfg *config.Config
}

// NewServiceContext .
func NewServiceContext(s store.Storage) *ServiceContext {
	c := config.GetConfig()
	return &ServiceContext{Rs: s.GetRedis(), Cfg: c}
}
