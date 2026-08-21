package sysuser

import (
	"context"
	"time"

	"gin-layout/internal/menu"
)

// Repository 查询接口
type Repository interface {
	Create(ctx context.Context, entity *SysUser) error
	Update(ctx context.Context, entity *SysUser) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*SysUser, bool, error)
	List(ctx context.Context, q userListQuery) ([]SysUser, int64, error)
	FindByAccount(ctx context.Context, account string) (*SysUser, bool, error)
	FindByIDWithRoles(ctx context.Context, id int64) (*SysUser, bool, error)
	UpdateLastLogin(ctx context.Context, userID int64, lastLoginAt time.Time) error
	CreateWithRoles(ctx context.Context, u *SysUser, roleIDs []int64) error
	UpdateWithRoles(ctx context.Context, u *SysUser, roleIDs []int64) error
	ReplaceUserRoles(ctx context.Context, userID int64, roleIDs []int64) error
	FindByIDs(ctx context.Context, ids []int64) ([]SysUser, error)
}

type RoleFinder interface {
	ListEnabledRoleIDsForUser(ctx context.Context, userID int64) ([]int64, error)
}

type MenuFinder interface {
	ListEnabledByRoleIDs(ctx context.Context, roleIDs []int64) ([]menu.Menu, error)
	ListAll(ctx context.Context) ([]menu.Menu, error)
	ToMenuTree(rows []menu.Menu) []menu.MenuItem
}
