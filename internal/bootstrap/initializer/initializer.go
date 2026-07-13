package initializer

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"

	"gin-layout/config"
	"gin-layout/internal/common"
	"gin-layout/internal/domain"
	"gin-layout/internal/infra"
	"gin-layout/internal/user"
	"gin-layout/internal/role"
	"gin-layout/internal/menu"
)

type MenuDefinition struct {
	ID         int64           `json:"id"`
	ParentID   *int64          `json:"parentId"`
	Name       string          `json:"name"`
	Code       *string         `json:"code"`
	Type       domain.MenuType `json:"type"`
	Path       string          `json:"path"`
	Redirect   string          `json:"redirect"`
	Component  string          `json:"component"`
	Icon       string          `json:"icon"`
	ActiveMenu string          `json:"activeMenu"`
	Link       string          `json:"link"`
	Query      string          `json:"query"`
	Remark     string          `json:"remark"`
	Sort       int             `json:"sort"`
	Level      int             `json:"level"`
	Hidden     bool            `json:"hidden"`
	Cache      bool            `json:"cache"`
	Affix      bool            `json:"affix"`
	Breadcrumb bool            `json:"breadcrumb"`
	AlwaysShow bool            `json:"alwaysShow"`
	External   bool            `json:"external"`
	Iframe     bool            `json:"iframe"`
	Enabled    bool            `json:"enabled"`
	Method     string          `json:"method"`
	APIPath    string          `json:"apiPath"`
	PermCode   *string         `json:"permCode"`
}

type Initializer struct {
	cfg            *config.Config
	userRepo       *user.PGRepository
	roleRepo       *role.PGRepository
	menuRepo       *menu.PGRepository
	passwordHasher common.PasswordHasher
	policies       common.PolicyManager
	logger         *infra.Logger
}

func NewInitializer(
	cfg *config.Config,
	userRepo *user.PGRepository,
	roleRepo *role.PGRepository,
	menuRepo *menu.PGRepository,
	passwordHasher common.PasswordHasher,
	policies common.PolicyManager,
	logger *infra.Logger,
) *Initializer {
	return &Initializer{
		cfg:            cfg,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		menuRepo:       menuRepo,
		passwordHasher: passwordHasher,
		policies:       policies,
		logger:         logger,
	}
}

func (i *Initializer) Run(ctx context.Context) error {
	if err := i.initSuperAdmin(ctx); err != nil {
		return fmt.Errorf("initialize super admin: %w", err)
	}
	if err := i.initMenu(ctx); err != nil {
		return fmt.Errorf("initialize menu: %w", err)
	}
	return nil
}

func (i *Initializer) initSuperAdmin(ctx context.Context) error {
	account := cmp.Or(i.cfg.Initializer.AdminAccount, "admin")
	password := cmp.Or(i.cfg.Initializer.AdminPassword, "admin123")
	roleName := cmp.Or(i.cfg.Initializer.AdminRoleName, "超级管理员")
	roleCode := cmp.Or(i.cfg.Initializer.AdminRoleCode, "ADMIN")

	r, err := i.roleRepo.FindByCode(ctx, roleCode)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	if errors.Is(err, domain.ErrNotFound) {
		r = &domain.Role{Name: roleName, Code: roleCode}
		if err := i.roleRepo.Create(ctx, r); err != nil {
			return err
		}
		i.logger.Info().Str("role", roleCode).Msg("created super admin role")
	} else {
		i.logger.Info().Str("role", roleCode).Msg("super admin role already exists")
	}

	u, err := i.userRepo.FindByAccount(ctx, account)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	if errors.Is(err, domain.ErrNotFound) {
		hash, err := i.passwordHasher.Hash(password)
		if err != nil {
			return err
		}
		u = &domain.User{
			Account:      account,
			PasswordHash: hash,
			NickName:     "超级管理员",
			Enabled:      true,
		}
		if err := i.userRepo.Create(ctx, u); err != nil {
			return err
		}
		i.logger.Info().Str("account", account).Msg("created super admin user")
	} else {
		i.logger.Info().Str("account", account).Msg("super admin user already exists")
	}

	if err := i.userRepo.ReplaceUserRoles(ctx, u.ID, []int64{r.ID}); err != nil {
		return fmt.Errorf("assign super admin role: %w", err)
	}

	if err := i.policies.SyncUserRoles(ctx, u.Account, []string{r.Code}); err != nil {
		return fmt.Errorf("sync super admin role policy: %w", err)
	}
	return nil
}

func (i *Initializer) initMenu(ctx context.Context) error {
	menuFile := cmp.Or(i.cfg.Initializer.MenuJSON, "internal/bootstrap/initializer/menu.json")

	menus, err := i.loadMenus(menuFile)
	if err != nil {
		return err
	}

	if len(menus) == 0 {
		i.logger.Info().Str("file", menuFile).Msg("no permission definitions found")
		return nil
	}

	data := lo.Map(menus, func(m MenuDefinition, _ int) domain.Menu {
		return domain.Menu{
			ID:         m.ID,
			ParentID:   m.ParentID,
			Name:       m.Name,
			Code:       optionalSeedString(m.Code),
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
			PermCode:   optionalSeedString(m.PermCode),
		}
	})

	if err := i.menuRepo.CreateAll(ctx, data); err != nil {
		return fmt.Errorf("create menus: %w", err)
	}

	i.logger.Info().Int("count", len(menus)).Str("file", menuFile).Msg("menu definitions loaded")
	return nil
}

func (i *Initializer) loadMenus(filePath string) ([]MenuDefinition, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read menu json file: %w", err)
	}

	var menus []MenuDefinition
	if err := json.Unmarshal(content, &menus); err != nil {
		return nil, fmt.Errorf("parse menu json file: %w", err)
	}
	return menus, nil
}

func optionalSeedString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
