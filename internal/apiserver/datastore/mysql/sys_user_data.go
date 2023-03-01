package mysql

import (
	"context"

	"github.com/BooeZhang/gin-layout/internal/apiserver/model"
	"gorm.io/gorm"
)

type sysUserDB struct {
	db *gorm.DB
}

func newSysUser(ds *_datastore) *sysUserDB {
	return &sysUserDB{db: ds.db}
}

// Create 创建用户
func (u *sysUserDB) Create(ctx context.Context, user *model.SysUserModel) error {

	return u.db.Create(&user).Error
}

// Update 更新用户
func (u *sysUserDB) Update(ctx context.Context, user *model.SysUserModel) error {

	return u.db.Save(user).Error
}

// GetSysUserByName 根据用户名获取用户
func (u *sysUserDB) GetSysUserByName(ctx context.Context, name string) (user *model.SysUserModel, err error) {
	err = u.db.Model(new(model.SysUserModel)).Where("user_name=?", name).First(&user).Error

	return
}
