package role

import (
	"context"

	"gin-layout/internal/domain"
)

type Repository interface {
	List(ctx context.Context, q roleListQuery) ([]domain.Role, int64, error)
	Create(ctx context.Context, role *domain.Role) error
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*domain.Role, error)
	FindByIDs(ctx context.Context, ids []int64) ([]domain.Role, error)
	CreateWithMenu(ctx context.Context, role *domain.Role, menuIDs []int64) error
	UpdateWithMenu(ctx context.Context, role *domain.Role, menuIDs []int64) error
	DeleteWithRelat(ctx context.Context, roleID int64) error
	FindByCode(ctx context.Context, code string) (*domain.Role, error)
	FindCodesByIDs(ctx context.Context, roleIDs []int64) ([]string, error)
	FindByUserIDs(ctx context.Context, userIDs []int64, enabled *bool) ([]domain.Role, error)
	FindByIDWithPerm(ctx context.Context, roleID int64) (*domain.Role, error)
	ListAll(ctx context.Context) ([]domain.Role, error)
	RoleAddUser(ctx context.Context, data []domain.UserRole) error
	RoleRemoveUser(ctx context.Context, roleID int64, userIDs []int64) error
	ReplaceRoleMenus(ctx context.Context, roleID int64, menuIDs []int64) error
}

type UserFinder interface {
	FindByIDs(ctx context.Context, ids []int64) ([]domain.User, error)
}

type MenuService interface {
	ListAll(ctx context.Context) ([]domain.Menu, error)
	FindAllMenuIDs(ctx context.Context) ([]int64, error)
	CompleteMenuIDsWithAncestors(ctx context.Context, menuIDs []int64) ([]int64, error)
	FindPermissionObjectsByMenuIDs(ctx context.Context, menuIDs []int64) ([]string, error)
}
