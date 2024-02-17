package mysql

import (
	"context"
	"github.com/BooeZhang/gin-layout/internal/model"
	"github.com/BooeZhang/gin-layout/internal/repo"
	"github.com/BooeZhang/gin-layout/store"
)

var _ repo.UserRepo = (*userRepo)(nil)

type userRepo struct {
	storeIns store.Storage
}

func NewUserRepo(s store.Storage) *userRepo {
	return &userRepo{
		storeIns: s,
	}
}

func (ur *userRepo) GetUserByName(ctx context.Context, name string) (*model.SysUser, error) {
	user := &model.SysUser{}
	err := ur.storeIns.GetMySQL().Where("name = ?", name).Find(user).Error
	return user, err
}

func (ur *userRepo) GetUserById(ctx context.Context, uid int64) (*model.SysUser, error) {
	user := &model.SysUser{}
	err := ur.storeIns.GetMySQL().Where("id = ?", uid).Find(user).Error
	return user, err
}

func (ur *userRepo) GetUserByMobile(ctx context.Context, mobile string) (*model.SysUser, error) {
	user := &model.SysUser{}
	err := ur.storeIns.GetMySQL().
		Where("mobile = ?", mobile).
		Where("enabled_status = 1").
		First(user).Error
	return user, err
}
