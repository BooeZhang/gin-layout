package datainterface

import (
	"context"
	"github.com/BooeZhang/gin-layout/internal/apiserver/model"
)

type SysUserData interface {
	Create(ctx context.Context, user *model.SysUserModel) error
	Update(ctx context.Context, user *model.SysUserModel) error
	GetSysUserByName(ctx context.Context, name string) (user *model.SysUserModel, err error)
}
