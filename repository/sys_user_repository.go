package repository

import (
	"context"
	"github.com/BooeZhang/gin-layout/model"
	"gorm.io/gorm"
)

type sysUserRepository struct {
	db *gorm.DB
}

func NewSysUserRepository(db *gorm.DB) model.SysUserRepository {
	return &sysUserRepository{
		db: db,
	}
}

func (su *sysUserRepository) Create(c context.Context, sysUser *model.SysUser) error {
	return su.db.Create(sysUser).Error
}

func (su *sysUserRepository) GetSysUserByName(c context.Context, name string) (*model.SysUser, error) {
	var d model.SysUser
	err := su.db.Model(new(model.SysUser)).Where("user_name=?", name).First(&d).Error
	return &d, err

}

// Update 更新用户
func (su *sysUserRepository) Update(ctx context.Context, user *model.SysUser) error {

	return su.db.Save(user).Error
}
