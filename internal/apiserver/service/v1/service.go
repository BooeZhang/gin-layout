package v1

import (
	"github.com/BooeZhang/gin-layout/internal/apiserver/datastore"
	"github.com/BooeZhang/gin-layout/internal/pkg/config"
)

// Service 资源接口
type Service interface {
	SysUser() SysUserSrv
}

type serviceContext struct {
	store datastore.Factory
	cfg   *config.Config
}

// NewService returns Service interface.
func NewService(store datastore.Factory, c *config.Config) Service {
	if c == nil {
		c = config.GetConfig()
	}
	return &serviceContext{store: store, cfg: c}
}

// SysUser 用户service
func (srv *serviceContext) SysUser() SysUserSrv {
	return newSysUserService(srv)
}
