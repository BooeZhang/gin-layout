package admin

import (
	"context"

	"github.com/BooeZhang/gin-layout/internal/model"
	"github.com/BooeZhang/gin-layout/internal/repo"
	"github.com/BooeZhang/gin-layout/internal/repo/mysql"
	"github.com/BooeZhang/gin-layout/internal/service/v1"
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/BooeZhang/gin-layout/store"
)

// UserService 用户业务逻辑
type UserService struct {
	ctx *v1.ServiceContext
	ur  repo.UserRepo
}

func NewUserService(s store.Storage) *UserService {
	return &UserService{
		ctx: v1.NewServiceContext(),
		ur:  mysql.NewUserRepo(s),
	}
}

// GetByName 通过用户名 查找用户
func (us *UserService) GetByName(ctx context.Context, name string) (*model.SysUser, error) {
	if len(name) == 0 {
		return nil, erroron.ErrNotFound
	}
	return us.ur.GetUserByName(ctx, name)
}

// GetById 根据用户ID查找用户
func (us *UserService) GetById(ctx context.Context, uid int64) (*model.SysUser, error) {
	return us.ur.GetUserById(ctx, uid)
}

// GetByMobile 根据用户手机号查询
func (us *UserService) GetByMobile(ctx context.Context, mobile string) (*model.SysUser, error) {
	// 认为handler层对service层入参都是合法的，除了业务上的校验，service层不校验入参合规性
	return us.ur.GetUserByMobile(ctx, mobile)
}
