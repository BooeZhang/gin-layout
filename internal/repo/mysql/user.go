package mysql

import (
	"context"

	"github.com/BooeZhang/gin-layout/internal/model"
	"github.com/BooeZhang/gin-layout/internal/repo"
	"github.com/BooeZhang/gin-layout/store"
)

var _ repo.UserRepo = (*UserRepo)(nil)

type UserRepo struct {
	storeIns store.Storage
}

func NewUserRepo(s store.Storage) *UserRepo {
	return &UserRepo{
		storeIns: s,
	}
}

func (ur *UserRepo) GetUserByName(ctx context.Context, name string) (*model.SysUser, error) {
	user := &model.SysUser{}
	err := ur.storeIns.GetMySQL().Where("name = ?", name).Find(user).Error
	return user, err
}

func (ur *UserRepo) GetUserById(ctx context.Context, uid int64) (*model.SysUser, error) {
	user := &model.SysUser{}
	err := ur.storeIns.GetMySQL().Where("id = ?", uid).Find(user).Error
	return user, err
}

func (ur *UserRepo) GetUserByMobile(ctx context.Context, mobile string) (*model.SysUser, error) {
	user := &model.SysUser{}
	err := ur.storeIns.GetMySQL().
		Where("mobile = ?", mobile).
		Where("enabled_status = 1").
		First(user).Error
	return user, err
}
