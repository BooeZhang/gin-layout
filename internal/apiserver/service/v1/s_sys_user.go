package v1

import (
	"context"

	"github.com/BooeZhang/gin-layout/internal/apiserver/model"
)

// SysUserSrv 用户资源接口
type SysUserSrv interface {
	Create(ctx context.Context, user *model.SysUserModel) error
	Update(ctx context.Context, user *model.SysUserModel) error
}

type sysUserService struct {
	*serviceContext
}

var _ SysUserSrv = (*sysUserService)(nil)

func newSysUserService(srv *serviceContext) *sysUserService {
	return &sysUserService{srv}
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
