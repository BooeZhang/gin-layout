package user

import (
	"context"
	"github.com/BooeZhang/gin-layout/internal/model"
	"github.com/BooeZhang/gin-layout/internal/repo"
	v1 "github.com/BooeZhang/gin-layout/internal/service/v1"
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/BooeZhang/gin-layout/pkg/schema"
)

var _ Service = (*serviceImpl)(nil)

// Service 定义用户操作服务接口
type Service interface {
	Login(ctx context.Context, name, pwd string) (*schema.LoginRes, error)
	GetByName(ctx context.Context, name string) (*model.SysUser, error)
	GetById(ctx context.Context, uid int64) (*model.SysUser, error)
	GetByMobile(ctx context.Context, ID string) (*model.SysUser, error)
}

// serviceImpl 实现UserService接口
type serviceImpl struct {
	ctx *v1.ServiceContext
	ur  repo.UserRepo
}

func NewUserService(ctx *v1.ServiceContext, _ur repo.UserRepo) *serviceImpl {
	return &serviceImpl{
		ctx: ctx,
		ur:  _ur,
	}
}

// GetByName 通过用户名 查找用户
func (us *serviceImpl) GetByName(ctx context.Context, name string) (*model.SysUser, error) {
	if len(name) == 0 {
		return nil, erroron.ErrNotFound
	}
	return us.ur.GetUserByName(ctx, name)
}

// GetById 根据用户ID查找用户
func (us *serviceImpl) GetById(ctx context.Context, uid int64) (*model.SysUser, error) {
	return us.ur.GetUserById(ctx, uid)
}

// GetByMobile 根据用户手机号查询
func (us *serviceImpl) GetByMobile(ctx context.Context, mobile string) (*model.SysUser, error) {
	// 认为handler层对service层入参都是合法的，除了业务上的校验，service层不校验入参合规性
	return us.ur.GetUserByMobile(ctx, mobile)
}
