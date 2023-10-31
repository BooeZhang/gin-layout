package mysql

import (
	"context"
	"github.com/BooeZhang/gin-layout/internal/model"
	"github.com/BooeZhang/gin-layout/internal/repo"
	"gorm.io/gorm"
)

var _ repo.UserRepo = (*userRepo)(nil)

type userRepo struct {
	ds *gorm.DB
}

func NewUserRepo(_ds *gorm.DB) *userRepo {
	return &userRepo{
		ds: _ds,
	}
}

func (ur *userRepo) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	user := &model.User{}
	err := ur.ds.Where("name = ?", name).Find(user).Error
	return user, err
}

func (ur *userRepo) GetUserById(ctx context.Context, uid int64) (*model.User, error) {
	user := &model.User{}
	err := ur.ds.Where("id = ?", uid).Find(user).Error
	return user, err
}

func (ur *userRepo) GetUserByMobile(ctx context.Context, mobile string) (*model.User, error) {
	user := &model.User{}
	err := ur.ds.
		Where("mobile = ?", mobile).
		Where("enabled_status = 1").
		First(user).Error
	return user, err
}
