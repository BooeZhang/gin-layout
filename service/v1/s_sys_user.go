package v1

import (
	"context"
	"github.com/BooeZhang/gin-layout/model"
	"github.com/BooeZhang/gin-layout/repository"
)

// SysUserSrv 用户资源接口
type SysUserSrv interface {
	Create(ctx context.Context, user *model.SysUser) error
	Update(ctx context.Context, user *model.SysUser) error
}

type sysUserService struct {
	*serviceContext
	sysUser model.SysUserRepository
}

var _ SysUserSrv = (*sysUserService)(nil)

func newSysUserService(srv *serviceContext) *sysUserService {
	return &sysUserService{srv, repository.NewSysUserRepository(srv.ms)}
}

func (su *sysUserService) Create(ctx context.Context, user *model.SysUser) error {
	if err := su.sysUser.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (su *sysUserService) Update(ctx context.Context, secret *model.SysUser) error {
	// Save changed fields.
	if err := su.sysUser.Update(ctx, secret); err != nil {
		return err
	}

	return nil
}
