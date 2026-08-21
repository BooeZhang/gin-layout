package migration

import (
	"fmt"

	"gin-layout/internal/infra"
	"gin-layout/internal/menu"
	"gin-layout/internal/role"
	"gin-layout/internal/sysuser"
)

// Run 数据库迁移
func Run(db *infra.Database) error {
	if db == nil || db.DB == nil {
		return infra.ErrDbIsNil
	}

	if err := db.Migrate(
		&sysuser.SysUser{},
		&sysuser.SysUserRole{},
		&role.Role{},
		&role.RoleMenu{},
		&menu.Menu{},
		&infra.TokenBlacklistModel{},
	); err != nil {
		return fmt.Errorf("auto migrate application schema: %w", err)
	}

	return nil
}
