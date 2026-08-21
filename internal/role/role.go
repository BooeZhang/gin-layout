package role

import "time"

// Role 角色,同时作为领域对象与 GORM 持久化模型。
type Role struct {
	ID          int64     `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Name        string    `gorm:"size:50;not null"`
	Code        string    `gorm:"uniqueIndex;size:50;not null"`
	Description string    `gorm:"size:255"`
	Sort        int       `gorm:"default:0"`
	Enabled     bool      `gorm:"default:true"`
	MenuIDs     []int64   `gorm:"-"`
}

func (Role) TableName() string { return "roles" }

// RoleMenu 角色-菜单关联表。
type RoleMenu struct {
	RoleID int64 `gorm:"primaryKey"`
	MenuID int64 `gorm:"primaryKey;index"`
}

func (RoleMenu) TableName() string { return "role_menus" }
