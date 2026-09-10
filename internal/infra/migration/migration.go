package migration

import (
	"fmt"

	"gin-layout/internal/domain/menu"
	"gin-layout/internal/domain/role"
	"gin-layout/internal/domain/sysuser"
	"gin-layout/internal/infra"
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
