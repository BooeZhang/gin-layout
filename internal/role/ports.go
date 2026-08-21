package role

import (
	"context"

	"gin-layout/internal/menu"
	"gin-layout/internal/sysuser"
)

type Repository interface {
	List(ctx context.Context, q roleListQuery) ([]Role, int64, error)
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*Role, bool, error)
	FindByIDs(ctx context.Context, ids []int64) ([]Role, error)
	CreateWithMenu(ctx context.Context, role *Role, menuIDs []int64) error
	UpdateWithMenu(ctx context.Context, role *Role, menuIDs []int64) error
	DeleteWithRoleID(ctx context.Context, roleID int64) error
	FindByCode(ctx context.Context, code string) (*Role, bool, error)
	FindCodesByIDs(ctx context.Context, roleIDs []int64) ([]string, error)
	FindByUserIDs(ctx context.Context, userIDs []int64, enabled *bool) ([]Role, error)
	FindByIDWithPerm(ctx context.Context, roleID int64) (*Role, bool, error)
	ListAll(ctx context.Context) ([]Role, error)
	RoleAddUser(ctx context.Context, data []sysuser.SysUserRole) error
	RoleRemoveUser(ctx context.Context, roleID int64, userIDs []int64) error
	ReplaceRoleMenus(ctx context.Context, roleID int64, menuIDs []int64) error
}

type UserFinder interface {
	FindByIDs(ctx context.Context, ids []int64) ([]sysuser.SysUser, error)
}

type MenuService interface {
	ListAll(ctx context.Context) ([]menu.Menu, error)
	FindAllMenuIDs(ctx context.Context) ([]int64, error)
	CompleteMenuIDsWithAncestors(ctx context.Context, menuIDs []int64) ([]int64, error)
	FindPermissionObjectsByMenuIDs(ctx context.Context, menuIDs []int64) ([]string, error)
}
