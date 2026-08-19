package sysuser

import (
	"context"
	"time"

	"gin-layout/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, entity *domain.SysUser) error
	Update(ctx context.Context, entity *domain.SysUser) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*domain.SysUser, error)
	List(ctx context.Context, q userListQuery) ([]domain.SysUser, int64, error)
	FindByAccount(ctx context.Context, account string) (*domain.SysUser, error)
	FindByIDWithRoles(ctx context.Context, id int64) (*domain.SysUser, error)
	UpdateLastLogin(ctx context.Context, userID int64, lastLoginAt time.Time) error
	CreateWithRoles(ctx context.Context, u *domain.SysUser, roleIDs []int64) error
	UpdateWithRoles(ctx context.Context, u *domain.SysUser, roleIDs []int64) error
	ReplaceUserRoles(ctx context.Context, userID int64, roleIDs []int64) error
	FindByIDs(ctx context.Context, ids []int64) ([]domain.SysUser, error)
}

type RoleFinder interface {
	ListEnabledRoleIDsForUser(ctx context.Context, userID int64) ([]int64, error)
}

type MenuFinder interface {
	ListEnabledByRoleIDs(ctx context.Context, roleIDs []int64) ([]domain.Menu, error)
	ListAll(ctx context.Context) ([]domain.Menu, error)
	ToMenuTree(rows []domain.Menu) []domain.MenuItem
}
