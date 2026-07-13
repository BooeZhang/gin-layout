package user

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

	userRepo Repository
	policy   common.PolicyManager
	password common.PasswordHasher
	roles    RoleFinder
	menus    MenuFinder
}

func NewService(base common.BaseService, users Repository, policy common.PolicyManager, password common.PasswordHasher, roles RoleFinder, menus MenuFinder) *Service {
	return &Service{
		BaseService: base,
		userRepo:    users,
		policy:      policy,
		password:    password,
		roles:       roles,
		menus:       menus,
	}
}

func (s *Service) List(ctx context.Context, in ListUserReq) (res domain.PageResult[UserItem], err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("user list")

	q := userListQuery{
		PageRequest: domain.PageRequest{Page: in.Page, PageSize: in.PageSize},
		Account:     in.Account,
		NickName:    in.NickName,
		Email:       in.Email,
		Phone:       in.Phone,
		Enabled:     in.Enabled,
	}

	result, total, err := s.userRepo.List(ctx, q)
	if err != nil {
		return res, err
	}

	items := lo.Map(result, func(item domain.User, _ int) UserItem {
		return s.toUserItem(item)
	})

	return domain.NewPageResult(items, total, in.Page, in.PageSize), nil
}

func (s *Service) Create(ctx context.Context, in CreateUserReq) (res CreateUserRes, err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("creating user")

	if existing, err := s.userRepo.FindByAccount(ctx, in.Account); err == nil && existing != nil {
		return res, domain.ErrAccountExists
	}

	hash, err := s.password.Hash(in.Password)
	if err != nil {
		return res, fmt.Errorf("hash password for %s: %w", in.Account, err)
	}

	u := &domain.User{
		Account:      in.Account,
		PasswordHash: hash,
		NickName:     in.NickName,
		Email:        in.Email,
		Phone:        in.Phone,
		Enabled:      true,
	}
	if len(u.NickName) == 0 {
		u.NickName = in.Account
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return res, err
	}

	if len(in.RoleIDs) > 0 {
		if err := s.userRepo.CreateWithRoles(ctx, u, in.RoleIDs); err != nil {
			return res, err
		}
		if err := s.policy.SyncUserRolesByIDs(ctx, u.Account, in.RoleIDs); err != nil {
			return res, fmt.Errorf("sync user roles (account=%s): %w", u.Account, err)
		}
	}

	return CreateUserRes{ID: u.ID}, nil
}

func (s *Service) GetDetails(ctx context.Context) (res UserItem, err error) {
	currUser, ok := domain.CurrentUserFromContext(ctx)
	if !ok {
		return res, domain.ErrNotLogin
	}
	user, err := s.userRepo.FindByIDWithRoles(ctx, currUser.UserID)
	if err != nil {
		return res, err
	}
	return s.toUserItem(*user), nil
}

func (s *Service) Update(ctx context.Context, in UpdateUserReq) (res UpdateUserRes, err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("update user")

	current, err := s.userRepo.FindByID(ctx, in.UserID)
	if err != nil {
		return res, err
	}

	if in.NickName != nil {
		current.NickName = *in.NickName
	}
	if in.Password != nil {
		pwd, err := s.password.Hash(*in.Password)
		if err != nil {
			return res, err
		}
		current.PasswordHash = pwd
	}
	if in.Email != nil {
		current.Email = *in.Email
	}
	if in.Phone != nil {
		current.Phone = *in.Phone
	}
	if in.Avatar != nil {
		current.Avatar = *in.Avatar
	}
	if in.Enabled != nil {
		current.Enabled = *in.Enabled
	}
	if len(in.RoleIDs) > 0 {
		if err := s.userRepo.UpdateWithRoles(ctx, current, in.RoleIDs); err != nil {
			return res, err
		}
		if err := s.policy.SyncUserRolesByIDs(ctx, current.Account, in.RoleIDs); err != nil {
			return res, fmt.Errorf("sync user roles (account=%s): %w", current.Account, err)
		}
	} else {
		if err := s.userRepo.Update(ctx, current); err != nil {
			return res, err
		}
	}
	return res, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	logger := s.Log(ctx)
	logger.Debug().Int64("id", id).Msg("deleting user")

	u, err := s.userRepo.FindByID(ctx, id)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	if errors.Is(err, domain.ErrNotFound) {
		return nil
	}

	if s.IsAdmin(u.Account) {
		return domain.ErrCannotDeleteAdmin
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("DeleteUser (id=%d): %w", id, err)
	}
	return nil
}

func (s *Service) GetCurrentUserMenus(ctx context.Context) ([]domain.MenuItem, error) {
	currUser, ok := domain.CurrentUserFromContext(ctx)
	if !ok {
		return nil, domain.ErrNotLogin
	}

	if s.IsAdmin(currUser.Account) {
		rows, err := s.menus.ListAll(ctx)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
		return s.menus.ToMenuTree(rows), nil
	}

	roleIDs, err := s.roles.ListEnabledRoleIDsForUser(ctx, currUser.UserID)
	if err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		return []domain.MenuItem{}, nil
	}

	rows, err := s.menus.ListEnabledByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}
	return s.menus.ToMenuTree(rows), nil
}

func (s *Service) toUserItem(u domain.User) UserItem {
	return UserItem{
		ID:          u.ID,
		Account:     u.Account,
		NickName:    u.NickName,
		Email:       u.Email,
		Phone:       u.Phone,
		Avatar:      u.Avatar,
		Enabled:     u.Enabled,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Roles: lo.Map(u.RoleIDs, func(roleID int64, _ int) RoleItem {
			return RoleItem{ID: roleID}
		}),
	}
}
