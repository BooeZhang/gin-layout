package mysql

import (
	"context"

	"github.com/BooeZhang/gin-layout/internal/apiserver/model"
	"gorm.io/gorm"
)

type ISysUser interface {
	Create(ctx context.Context, user *model.SysUserModel) error
	Update(ctx context.Context, user *model.SysUserModel) error
	GetSysUserByName(ctx context.Context, name string) (user *model.SysUserModel, err error)
}

type sysUser struct {
	db *gorm.DB
}

func newSysUser(ds *datastore) *sysUser {
	return &sysUser{db: ds.db}
}

// Create 创建用户
func (u *sysUser) Create(ctx context.Context, user *model.SysUserModel) error {

	return u.db.Create(&user).Error
}

// Update 更新用户
func (u *sysUser) Update(ctx context.Context, user *model.SysUserModel) error {

	return u.db.Save(user).Error
}

// GetSysUserByName 根据用户名获取用户
func (u *sysUser) GetSysUserByName(ctx context.Context, name string) (user *model.SysUserModel, err error) {
	err = u.db.Model(new(model.SysUserModel)).Where("user_name=?", name).First(&user).Error

	return
}
