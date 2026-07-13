package menu

import (
	"context"

	"gin-layout/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, m *domain.Menu) error
	Update(ctx context.Context, m *domain.Menu) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*domain.Menu, error)
	FindByIDs(ctx context.Context, ids []int64) ([]domain.Menu, error)
	ListAll(ctx context.Context) ([]domain.Menu, error)
	FindByCode(ctx context.Context, code string) (*domain.Menu, error)
	FindMenusByRoleIDs(ctx context.Context, roleIDs []int64, enabled *bool) ([]domain.Menu, error)
	CreateAll(ctx context.Context, menus []domain.Menu) error
}
