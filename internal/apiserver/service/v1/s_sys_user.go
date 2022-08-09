package v1

import (
	"context"

	"github.com/BooeZhang/gin-layout/internal/apiserver/datastore"
	"github.com/BooeZhang/gin-layout/internal/apiserver/model"
)

// SysUserv 用户资源接口
type SysUserv interface {
	Create(ctx context.Context, user *model.SysUserModel) error
	Update(ctx context.Context, user *model.SysUserModel) error
}

type sysUserService struct {
	store datastore.Factory
}

var _ SysUserv = (*sysUserService)(nil)

func newSysUserService(srv *service) *sysUserService {
	return &sysUserService{store: srv.store}
}

func (su *sysUserService) Create(ctx context.Context, user *model.SysUserModel) error {
	if err := su.store.SysUser().Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (su *sysUserService) Update(ctx context.Context, secret *model.SysUserModel) error {
	// Save changed fields.
	if err := su.store.SysUser().Update(ctx, secret); err != nil {
		return err
	}

	return nil
}
