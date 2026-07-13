package user

import (
	"context"
	"time"

	"gin-layout/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, entity *domain.User) error
	Update(ctx context.Context, entity *domain.User) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	List(ctx context.Context, q userListQuery) ([]domain.User, int64, error)
	FindByAccount(ctx context.Context, account string) (*domain.User, error)
	FindByIDWithRoles(ctx context.Context, id int64) (*domain.User, error)
	UpdateLastLogin(ctx context.Context, userID int64, lastLoginAt time.Time) error
	CreateWithRoles(ctx context.Context, u *domain.User, roleIDs []int64) error
	UpdateWithRoles(ctx context.Context, u *domain.User, roleIDs []int64) error
	ReplaceUserRoles(ctx context.Context, userID int64, roleIDs []int64) error
	FindByIDs(ctx context.Context, ids []int64) ([]domain.User, error)
}

type RoleFinder interface {
	ListEnabledRoleIDsForUser(ctx context.Context, userID int64) ([]int64, error)
}

type MenuFinder interface {
	ListEnabledByRoleIDs(ctx context.Context, roleIDs []int64) ([]domain.Menu, error)
	ListAll(ctx context.Context) ([]domain.Menu, error)
	ToMenuTree(rows []domain.Menu) []domain.MenuItem
}
