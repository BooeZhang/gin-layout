package repo

import (
	"context"
	"github.com/BooeZhang/gin-layout/internal/model"
)

// UserRepo 用户repo接口
type UserRepo interface {
	GetUserByName(ctx context.Context, name string) (*model.SysUser, error)
	GetUserById(ctx context.Context, uid int64) (*model.SysUser, error)
	GetUserByMobile(ctx context.Context, mobile string) (*model.SysUser, error)
}
