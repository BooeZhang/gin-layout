package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	"gin-layout/internal/common"
	"gin-layout/internal/domain"
)

type Service struct {
	common.BaseService

	repo      Repository
	userQuery UserFinder
	menus     MenuService
	policies  common.PolicyManager
}

func NewService(base common.BaseService, repo Repository, userQuery UserFinder, menus MenuService, policies common.PolicyManager) *Service {
	return &Service{BaseService: base, repo: repo, userQuery: userQuery, menus: menus, policies: policies}
}

func (s *Service) List(ctx context.Context, in ListRoleReq) (res domain.PageResult[RoleItem], err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("role list")

	q := roleListQuery{
		PageRequest: domain.PageRequest{Page: in.Page, PageSize: in.PageSize},
		Name:        in.Name, Code: in.Code, Enabled: in.Enabled,
	}
	result, total, err := s.repo.List(ctx, q)
	if err != nil {
		return res, err
	}

	items := lo.Map(result, func(item domain.Role, _ int) RoleItem {
		return RoleItem{
			ID: item.ID, Name: item.Name, Code: item.Code,
			Description: item.Description, Sort: item.Sort,
			Enabled: item.Enabled, CreatedAt: item.CreatedAt,
		}
	})
	return domain.NewPageResult(items, total, in.Page, in.PageSize), nil
}

func (s *Service) Create(ctx context.Context, in CreateRoleReq) (res CreateRoleRes, err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("creating role")

	_, err = s.repo.FindByCode(ctx, in.Code)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return res, err
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return res, domain.ErrRoleExists
	}

	role := &domain.Role{
		Name: in.Name, Code: in.Code, Description: in.Description,
		Sort: in.Sort, Enabled: in.Enabled,
	}

	if in.PermissionIDs != nil {
		if err := s.repo.CreateWithMenu(ctx, role, in.PermissionIDs); err != nil {
			return res, err
		}
		if err := s.syncRoleMenus(ctx, role.ID, role.Code, in.PermissionIDs); err != nil {
			return res, err
		}
	} else {
		if err := s.repo.Create(ctx, role); err != nil {
			return res, fmt.Errorf("CreateRole (code=%s): %w", in.Code, err)
		}
	}
	return CreateRoleRes{ID: role.ID}, nil
}

func (s *Service) GetOne(ctx context.Context, id int64) (*RoleItem, error) {
	logger := s.Log(ctx)
	logger.Debug().Int64("id", id).Msg("getting role")

	r, err := s.repo.FindByIDWithPerm(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &RoleItem{
		ID: r.ID, Name: r.Name, Code: r.Code,
		Description: r.Description, Sort: r.Sort, Enabled: r.Enabled,
		CreatedAt: r.CreatedAt, PermissionIDs: r.MenuIDs,
	}

	if s.IsAdminRole(res.Code) {
		allMenu, err := s.menus.ListAll(ctx)
		if err != nil {
			return nil, err
		}
		res.PermissionIDs = lo.Map(allMenu, func(item domain.Menu, _ int) int64 { return item.ID })
	}
	return res, nil
}

func (s *Service) Update(ctx context.Context, in UpdateRoleReq) (res UpdateRoleRes, err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("updating role")

	r, err := s.repo.FindByID(ctx, in.RoleID)
	if err != nil {
		return res, err
	}

	if in.Name != nil {
		r.Name = *in.Name
	}
	if in.Description != nil {
		r.Description = *in.Description
	}
	if in.Sort != nil {
		r.Sort = *in.Sort
	}
	if in.Enabled != nil {
		r.Enabled = *in.Enabled
	}

	if in.PermissionIDs != nil {
		if err := s.repo.UpdateWithMenu(ctx, r, in.PermissionIDs); err != nil {
			return res, err
		}
		if err := s.syncRoleMenus(ctx, r.ID, r.Code, in.PermissionIDs); err != nil {
			return res, err
		}
	} else {
		if err := s.repo.Update(ctx, r); err != nil {
			return res, err
		}
	}
	return res, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	logger := s.Log(ctx)
	logger.Debug().Int64("roleID", id).Msg("deleting role")

	role, err := s.repo.FindByID(ctx, id)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	if errors.Is(err, domain.ErrNotFound) {
		return nil
	}

	if s.IsAdminRole(role.Code) {
		return domain.ErrCannotDeleteAdminRole
	}

	if err := s.repo.DeleteWithRelat(ctx, id); err != nil {
		return fmt.Errorf("DeleteRole (id=%d): %w", id, err)
	}
	if err := s.policies.DeleteRole(ctx, role.Code); err != nil {
		return fmt.Errorf("DeleteRolePolicy (role=%s): %w", role.Code, err)
	}
	return nil
}

func (s *Service) GetAll(ctx context.Context) ([]RoleItem, error) {
	roles, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	items := lo.Map(roles, func(item domain.Role, _ int) RoleItem {
		return RoleItem{ID: item.ID, Name: item.Name, Code: item.Code}
	})
	return items, nil
}

func (s *Service) ListEnabledRoleIDsForUser(ctx context.Context, userID int64) ([]int64, error) {
	enabled := true
	roles, err := s.repo.FindByUserIDs(ctx, []int64{userID}, &enabled)
	if err != nil {
		return nil, err
	}
	return lo.Map(roles, func(item domain.Role, _ int) int64 { return item.ID }), nil
}

func (s *Service) UserAdd(ctx context.Context, in UserAddReq) (res UserAddRes, err error) {
	logger := s.Log(ctx)
	logger.Debug().Ints64("userIDs", in.UserIDs).Int64("roleID", in.RoleID).Msg("adding role user")

	role, err := s.repo.FindByID(ctx, in.RoleID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return res, err
	}
	if errors.Is(err, domain.ErrNotFound) {
		return res, nil
	}

	items := lo.Map(in.UserIDs, func(userID int64, _ int) domain.UserRole {
		return domain.UserRole{UserID: userID, RoleID: in.RoleID}
	})
	if err := s.repo.RoleAddUser(ctx, items); err != nil {
		return res, fmt.Errorf("RoleAddUser (roleID=%d, userIDs=%v): %w", in.RoleID, in.UserIDs, err)
	}

	users, err := s.userQuery.FindByIDs(ctx, in.UserIDs)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return res, fmt.Errorf("FindUsers for RoleAdd (roleID=%d, userIDs=%v): %w", in.RoleID, in.UserIDs, err)
	}
	for _, user := range users {
		if err := s.policies.AddRoleToUser(ctx, user.Account, role.Code); err != nil {
			return res, fmt.Errorf("AddRoleToUser (account=%s, role=%s): %w", user.Account, role.Code, err)
		}
	}
	return res, nil
}

func (s *Service) UserRemove(ctx context.Context, in UserRemoveReq) (res UserRemoveRes, err error) {
	logger := s.Log(ctx)
	logger.Debug().Ints64("userIDs", in.UserIDs).Int64("roleID", in.RoleID).Msg("removing role user")

	role, err := s.repo.FindByID(ctx, in.RoleID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return res, err
	}
	if errors.Is(err, domain.ErrNotFound) {
		return res, nil
	}

	if err := s.repo.RoleRemoveUser(ctx, in.RoleID, in.UserIDs); err != nil {
		return res, fmt.Errorf("RoleRemoveUser (roleID=%d, userIDs=%v): %w", in.RoleID, in.UserIDs, err)
	}

	users, err := s.userQuery.FindByIDs(ctx, in.UserIDs)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return res, fmt.Errorf("FindUsers for RoleRemove (roleID=%d, userIDs=%v): %w", in.RoleID, in.UserIDs, err)
	}
	for _, user := range users {
		roles, err := s.repo.FindByUserIDs(ctx, []int64{user.ID}, nil)
		if err != nil {
			return res, err
		}
		roleCodes := lo.Map(roles, func(item domain.Role, _ int) string { return item.Code })
		roleCodes = lo.Filter(roleCodes, func(code string, _ int) bool { return code != role.Code })
		if err := s.policies.SyncUserRoles(ctx, user.Account, roleCodes); err != nil {
			return res, fmt.Errorf("SyncUserRoles (account=%s, role=%s): %w", user.Account, role.Code, err)
		}
	}
	return res, nil
}

func (s *Service) syncRoleMenus(ctx context.Context, roleID int64, roleCode string, menuIDs []int64) error {
	completedMenuIDs, err := s.menus.CompleteMenuIDsWithAncestors(ctx, menuIDs)
	if err != nil {
		return err
	}
	if err := s.repo.ReplaceRoleMenus(ctx, roleID, completedMenuIDs); err != nil {
		return err
	}
	objects, err := s.menus.FindPermissionObjectsByMenuIDs(ctx, completedMenuIDs)
	if err != nil {
		return err
	}
	permissions := lo.Map(objects, func(object string, _ int) []string {
		return []string{object, "exec"}
	})
	return s.policies.SyncRolePermissions(ctx, roleCode, permissions)
}
