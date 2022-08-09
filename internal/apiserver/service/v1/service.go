package v1

import "github.com/BooeZhang/gin-layout/internal/apiserver/datastore"

// Service 资源接口
type Service interface {
	SysUser() SysUserv
}

type service struct {
	store datastore.Factory
}

// NewService returns Service interface.
func NewService(store datastore.Factory) Service {
	return &service{store: store}
}

// SysUser 用户service
func (srv *service) SysUser() SysUserv {
	return newSysUserService(srv)
}
