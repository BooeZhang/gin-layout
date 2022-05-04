package mysql

import (
	"fmt"
	"github.com/BooeZhang/gin-layout/internal/apiserver/model"
	"gorm.io/gorm"
)

func migrateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(
		new(model.SysUserModel),
	); err != nil {
		return fmt.Errorf("migrate user model failed: %w", err)
	}
	return nil
}
