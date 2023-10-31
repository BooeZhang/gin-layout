package v1

import (
	"github.com/BooeZhang/gin-layout/config"
	"github.com/go-redis/redis/v8"
)

type ServiceContext struct {
	rs  redis.UniversalClient
	cfg *config.Config
}

// NewServiceContext .
func NewServiceContext(rs redis.UniversalClient) *ServiceContext {
	c := config.GetConfig()
	return &ServiceContext{rs: rs, cfg: c}
}
