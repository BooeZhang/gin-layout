package v1

import (
	"github.com/BooeZhang/gin-layout/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Service 资源接口
type Service interface {
	SysUser() SysUserSrv
}

type serviceContext struct {
	ms  *gorm.DB
	rs  redis.UniversalClient
	cfg *config.Config
}

// NewService returns Service interface.
func NewService(ms *gorm.DB, rs redis.UniversalClient, c *config.Config) Service {
	if c == nil {
		c = config.GetConfig()
	}
	return &serviceContext{ms: ms, rs: rs, cfg: c}
}

// SysUser 用户service
func (srv *serviceContext) SysUser() SysUserSrv {
	return newSysUserService(srv)
}
