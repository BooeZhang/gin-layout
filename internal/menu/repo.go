package menu

import (
	"context"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gin-layout/internal/domain"
	"gin-layout/internal/infra"
)

type MenuModel struct {
	ID         int64           `gorm:"primaryKey"`
	CreatedAt  time.Time       `gorm:"autoCreateTime"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime"`
	ParentID   *int64          `gorm:"index"`
	Name       string          `gorm:"size:100;index;not null"`
	Code       *string         `gorm:"size:100;uniqueIndex:uk_menu_code"`
	Type       domain.MenuType `gorm:"size:20;default:'menu';index"`
	Path       string          `gorm:"size:255;index"`
	Redirect   string          `gorm:"size:255"`
	Component  string          `gorm:"size:255"`
	Icon       string          `gorm:"size:100"`
	ActiveMenu string          `gorm:"column:active_menu;size:255"`
	Link       string          `gorm:"size:500"`
	Query      string          `gorm:"type:text"`
	Remark     string          `gorm:"size:255"`
	Sort       int             `gorm:"default:0;index"`
	Level      int             `gorm:"default:0"`
	Hidden     bool            `gorm:"default:false"`
	Cache      bool            `gorm:"default:true"`
	Affix      bool            `gorm:"default:false"`
	Breadcrumb bool            `gorm:"default:true"`
	AlwaysShow bool            `gorm:"column:always_show;default:false"`
	External   bool            `gorm:"default:false"`
	Iframe     bool            `gorm:"default:false"`
	Enabled    bool            `gorm:"default:true;index"`
	Method     string          `gorm:"size:10"`
	APIPath    string          `gorm:"size:255"`
	PermCode   *string         `gorm:"size:255;uniqueIndex:uk_menu_perm_code"`
}

func (MenuModel) TableName() string { return "menus" }

func (m MenuModel) toDomain() domain.Menu {
	return domain.Menu{
		ID:         m.ID,
		ParentID:   m.ParentID,
		Name:       m.Name,
		Code:       m.Code,
		Type:       m.Type,
		Path:       m.Path,
		Redirect:   m.Redirect,
		Component:  m.Component,
		Icon:       m.Icon,
		ActiveMenu: m.ActiveMenu,
		Link:       m.Link,
		Query:      m.Query,
		Remark:     m.Remark,
		Sort:       m.Sort,
		Level:      m.Level,
		Hidden:     m.Hidden,
		Cache:      m.Cache,
		Affix:      m.Affix,
		Breadcrumb: m.Breadcrumb,
		AlwaysShow: m.AlwaysShow,
		External:   m.External,
		Iframe:     m.Iframe,
		Enabled:    m.Enabled,
		Method:     m.Method,
		APIPath:    m.APIPath,
		PermCode:   m.PermCode,
	}
}

func toMenuModel(m domain.Menu) MenuModel {
	return MenuModel{
		ID:         m.ID,
		ParentID:   m.ParentID,
		Name:       m.Name,
		Code:       m.Code,
		Type:       m.Type,
		Path:       m.Path,
		Redirect:   m.Redirect,
		Component:  m.Component,
		Icon:       m.Icon,
		ActiveMenu: m.ActiveMenu,
		Link:       m.Link,
		Query:      m.Query,
		Remark:     m.Remark,
		Sort:       m.Sort,
		Level:      m.Level,
		Hidden:     m.Hidden,
		Cache:      m.Cache,
		Affix:      m.Affix,
		Breadcrumb: m.Breadcrumb,
		AlwaysShow: m.AlwaysShow,
		External:   m.External,
		Iframe:     m.Iframe,
		Enabled:    m.Enabled,
		Method:     m.Method,
		APIPath:    m.APIPath,
		PermCode:   m.PermCode,
	}
}

type PGRepository struct {
	crud *infra.CRUDRepository[MenuModel, int64]
	db   *gorm.DB
}

func NewRepository(db *gorm.DB) *PGRepository {
	return &PGRepository{
		db:   db,
		crud: infra.NewCRUDRepository[MenuModel, int64](db),
	}
}

func (r *PGRepository) Create(ctx context.Context, m *domain.Menu) error {
	model := toMenuModel(*m)
	if err := r.crud.Create(ctx, &model); err != nil {
		return err
	}
	m.ID = model.ID
	return nil
}

func (r *PGRepository) Update(ctx context.Context, m *domain.Menu) error {
	model := toMenuModel(*m)
	return r.crud.Update(ctx, &model)
}

func (r *PGRepository) Delete(ctx context.Context, id int64) error {
	return r.crud.Delete(ctx, id)
}

func (r *PGRepository) FindByID(ctx context.Context, id int64) (*domain.Menu, error) {
	model, err := r.crud.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	m := model.toDomain()
	return &m, nil
}

func (r *PGRepository) FindByIDs(ctx context.Context, ids []int64) ([]domain.Menu, error) {
	models, err := r.crud.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	menus := lo.Map(models, func(m MenuModel, _ int) domain.Menu { return m.toDomain() })
	return menus, nil
}

func (r *PGRepository) ListAll(ctx context.Context) ([]domain.Menu, error) {
	models, err := gorm.G[MenuModel](r.db).Order("sort asc, id asc").Find(ctx)
	if err != nil {
		return nil, infra.NormalizeError(err)
	}
	menus := lo.Map(models, func(m MenuModel, _ int) domain.Menu { return m.toDomain() })
	return menus, nil
}

func (r *PGRepository) FindByCode(ctx context.Context, code string) (*domain.Menu, error) {
	m, err := gorm.G[MenuModel](r.db).Where("code = ?", code).First(ctx)
	if err != nil {
		return nil, infra.NormalizeError(err)
	}
	menu := m.toDomain()
	return &menu, nil
}

func (r *PGRepository) FindMenusByRoleIDs(ctx context.Context, roleIDs []int64, enabled *bool) ([]domain.Menu, error) {
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

	var models []MenuModel
	err := db.Select("DISTINCT menus.*").Order("menus.sort ASC, menus.id ASC").Scan(&models).Error
	if err != nil {
		return nil, infra.NormalizeError(err)
	}

	menus := lo.Map(models, func(m MenuModel, _ int) domain.Menu { return m.toDomain() })
	return menus, nil
}

func (r *PGRepository) CreateAll(ctx context.Context, menus []domain.Menu) error {
	if len(menus) == 0 {
		return nil
	}
	models := lo.Map(menus, func(m domain.Menu, _ int) MenuModel { return toMenuModel(m) })
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "perm_code"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"parent_id", "name", "type", "path", "redirect", "component",
				"icon", "active_menu", "link", "query", "remark", "sort", "level",
				"hidden", "cache", "affix", "breadcrumb", "always_show", "external",
				"iframe", "enabled", "method", "api_path", "perm_code",
			}),
		}).
		Create(&models).Error
	return infra.NormalizeError(err)
}
