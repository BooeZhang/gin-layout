package repo

import (
	"context"
	"github.com/BooeZhang/gin-layout/internal/model"
)

// UserRepo 用户repo接口
type UserRepo interface {
	GetUserByName(ctx context.Context, name string) (*model.User, error)
	GetUserById(ctx context.Context, uid int64) (*model.User, error)
	GetUserByMobile(ctx context.Context, mobile string) (*model.User, error)
}
