package menu

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gin-layout/internal/infra"
)

type MenuRepository struct {
	*infra.CRUDRepository[Menu, int64]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{
		db:             db,
		CRUDRepository: infra.NewCRUDRepository[Menu, int64](db),
	}
}

func (r *MenuRepository) ListAll(ctx context.Context) ([]Menu, error) {
	return gorm.G[Menu](r.db).Order("sort asc, id asc").Find(ctx)
}

func (r *MenuRepository) FindByCode(ctx context.Context, code string) (*Menu, bool, error) {
	m, err := gorm.G[Menu](r.db).Where("code = ?", code).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &m, true, nil
}

func (r *MenuRepository) FindMenusByRoleIDs(ctx context.Context, roleIDs []int64, enabled *bool) ([]Menu, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}

	db := r.db.WithContext(ctx).
		Table("menus").
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Where("role_menus.role_id IN ?", roleIDs)

	if enabled != nil {
		db = db.Where("menus.enabled = ?", *enabled)
	}

	var models []Menu
	err := db.Select("DISTINCT menus.*").Order("menus.sort ASC, menus.id ASC").Scan(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (r *MenuRepository) CreateAll(ctx context.Context, menus []Menu) error {
	if len(menus) == 0 {
		return nil
	}
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "perm_code"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"parent_id", "name", "type", "path", "redirect", "component",
				"icon", "active_menu", "link", "query", "remark", "sort", "level",
				"hidden", "cache", "affix", "breadcrumb", "always_show", "external",
				"iframe", "enabled", "method", "api_path", "perm_code",
			}),
		}).
		Create(&menus).Error
	return err
}
