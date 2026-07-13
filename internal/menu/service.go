package menu

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"gin-layout/internal/common"
	"gin-layout/internal/domain"
)

type Service struct {
	common.BaseService

	repo        Repository
	permMapMu   sync.RWMutex
	permCodeMap map[string]string
}

func NewService(base common.BaseService, repo Repository) *Service {
	return &Service{BaseService: base, repo: repo}
}

func permissionMapKey(url, method string) string {
	return strings.ToUpper(strings.TrimSpace(method)) + " " + strings.TrimSpace(url)
}

func (s *Service) LoadPermissionMap(ctx context.Context) error {
	permissions, err := s.repo.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("LoadPermissionMap: %w", err)
	}

	permCodeMap := make(map[string]string, len(permissions))
	for _, menu := range permissions {
		if menu.PermCode != nil {
			permCodeMap[permissionMapKey(menu.APIPath, menu.Method)] = *menu.PermCode
		}
	}

	s.permMapMu.Lock()
	s.permCodeMap = permCodeMap
	s.permMapMu.Unlock()
	return nil
}

func (s *Service) ResolvePermissionCode(url, method string) (string, bool) {
	s.permMapMu.RLock()
	defer s.permMapMu.RUnlock()
	code, ok := s.permCodeMap[permissionMapKey(url, method)]
	return code, ok
}

func (s *Service) List(ctx context.Context) ([]domain.MenuItem, error) {
	logger := s.Log(ctx)
	logger.Debug().Msg("listing menus")
	rows, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return ToMenuTree(rows), nil
}

func (s *Service) ListAll(ctx context.Context) ([]domain.Menu, error) {
	rows, err := s.repo.ListAll(ctx)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	return rows, nil
}

func (s *Service) Create(ctx context.Context, in CreateMenuReq) (res CreateMenuRes, err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("creating menu")

	if existing, err := s.repo.FindByCode(ctx, *in.Code); err == nil && existing != nil {
		return res, domain.ErrMenuExists
	}

	m := &domain.Menu{}
	applyCreateInput(m, in)

	if err := s.repo.Create(ctx, m); err != nil {
		return res, fmt.Errorf("CreateMenu (name=%s): %w", in.Name, err)
	}
	if err := s.LoadPermissionMap(ctx); err != nil {
		return res, fmt.Errorf("refresh perm map after create: %w", err)
	}
	return CreateMenuRes{ID: m.ID}, nil
}

func (s *Service) GetOne(ctx context.Context, id int64) (*domain.MenuItem, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toMenuItem(m), nil
}

func (s *Service) Update(ctx context.Context, in UpdateMenuReq) (res UpdateMenuRes, err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("updating menu")

	m, err := s.repo.FindByID(ctx, in.MenuID)
	if err != nil {
		return res, err
	}
	applyUpdateFields(m, in)

	if err := s.repo.Update(ctx, m); err != nil {
		return res, err
	}
	if err := s.LoadPermissionMap(ctx); err != nil {
		return res, fmt.Errorf("refresh perm map after update: %w", err)
	}
	return res, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	logger := s.Log(ctx)
	logger.Debug().Int64("menuID", id).Msg("deleting menu")
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("DeleteMenu (id=%d): %w", id, err)
	}
	if err := s.LoadPermissionMap(ctx); err != nil {
		return fmt.Errorf("refresh perm map after delete: %w", err)
	}
	return nil
}

func (s *Service) FindAllMenuIDs(ctx context.Context) ([]int64, error) {
	rows, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	menuIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		menuIDs = append(menuIDs, row.ID)
	}
	return menuIDs, nil
}

func (s *Service) CompleteMenuIDsWithAncestors(ctx context.Context, menuIDs []int64) ([]int64, error) {
	if len(menuIDs) == 0 {
		return []int64{}, nil
	}

	rows, err := s.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	parentByID := make(map[int64]*int64, len(rows))
	for _, row := range rows {
		parentByID[row.ID] = row.ParentID
	}

	seen := make(map[int64]struct{}, len(menuIDs))
	completed := make([]int64, 0, len(menuIDs))
	var addWithAncestors func(int64)
	addWithAncestors = func(menuID int64) {
		if _, ok := seen[menuID]; ok {
			return
		}
		seen[menuID] = struct{}{}
		completed = append(completed, menuID)
		if parentID := parentByID[menuID]; parentID != nil {
			addWithAncestors(*parentID)
		}
	}

	for _, menuID := range menuIDs {
		addWithAncestors(menuID)
	}
	return completed, nil
}

func (s *Service) FindPermissionObjectsByMenuIDs(ctx context.Context, menuIDs []int64) ([]string, error) {
	if len(menuIDs) == 0 {
		return []string{}, nil
	}

	rows, err := s.repo.FindByIDs(ctx, menuIDs)
	if err != nil {
		return nil, err
	}

	objects := make([]string, 0, len(rows))
	seen := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		object := row.APIPath
		if row.PermCode != nil && *row.PermCode != "" {
			object = *row.PermCode
		}
		if object == "" {
			continue
		}
		if _, ok := seen[object]; ok {
			continue
		}
		seen[object] = struct{}{}
		objects = append(objects, object)
	}
	return objects, nil
}

func (s *Service) ListEnabledByRoleIDs(ctx context.Context, roleIDs []int64) ([]domain.Menu, error) {
	if len(roleIDs) == 0 {
		return []domain.Menu{}, nil
	}
	enabled := true
	return s.repo.FindMenusByRoleIDs(ctx, roleIDs, &enabled)
}

func (s *Service) FindByRoleIDs(ctx context.Context, roleIDs []int64, enabled *bool) ([]domain.Menu, error) {
	return s.repo.FindMenusByRoleIDs(ctx, roleIDs, enabled)
}

func (s *Service) ToMenuTree(rows []domain.Menu) []domain.MenuItem {
	return ToMenuTree(rows)
}
