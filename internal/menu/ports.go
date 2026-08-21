package menu

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, m *Menu) error
	Update(ctx context.Context, m *Menu) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*Menu, bool, error)
	FindByIDs(ctx context.Context, ids []int64) ([]Menu, error)
	ListAll(ctx context.Context) ([]Menu, error)
	FindByCode(ctx context.Context, code string) (*Menu, bool, error)
	FindMenusByRoleIDs(ctx context.Context, roleIDs []int64, enabled *bool) ([]Menu, error)
	CreateAll(ctx context.Context, menus []Menu) error
}
