package sysuser

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/samber/lo"

	"gin-layout/internal/common"
	"gin-layout/internal/domain"
)

type Service struct {
	common.BaseService

	sysUserRepo Repository
	policy      common.PolicyManager
	roles       RoleFinder
	menus       MenuFinder
}

func NewService(base common.BaseService, users Repository, policy common.PolicyManager, roles RoleFinder, menus MenuFinder) *Service {
	return &Service{
		BaseService: base,
		sysUserRepo: users,
		policy:      policy,
		roles:       roles,
		menus:       menus,
	}
}

func (s *Service) Login(ctx context.Context, req LoginReq) (*LoginRes, error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", req).Msg("login attempt")

	u, err := s.sysUserRepo.FindByAccount(ctx, req.Account)
	if err != nil {
		return nil, fmt.Errorf("login: %w", domain.ErrAccountOrPassword)
	}

	if !u.Enabled {
		return nil, domain.ErrUserDisabled
	}

	if !u.ComparePassword(req.Password) {
		return nil, domain.ErrAccountOrPassword
	}

	tokens, err := s.TokenManager.Issue(u.ID, u.Account)
	if err != nil {
		return nil, fmt.Errorf("issue token (userID=%d): %w", u.ID, err)
	}

	_ = s.sysUserRepo.UpdateLastLogin(ctx, u.ID, time.Now())

	return &LoginRes{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *Service) Logout(ctx context.Context, _ LogoutReq) (*LogoutRes, error) {
	if _, err := s.TokenManager.RevokeCurrent(ctx); err != nil {
		return nil, err
	}
	return &LogoutRes{}, nil
}

func (s *Service) RefreshToken(ctx context.Context, req RefreshTokenReq) (*RefreshTokenRes, error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", req).Msg("refresh token")

	claims, err := s.TokenManager.Parse(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if claims.Type != domain.TokenTypeRefresh {
		return nil, domain.ErrTokenInvalid
	}

	revoked, err := s.TokenManager.IsRevoked(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if revoked {
		return nil, domain.ErrTokenRevoked
	}

	if err := s.TokenManager.Revoke(ctx, req.RefreshToken, claims.UserID, claims.ExpiresAt); err != nil {
		return nil, err
	}

	tokens, err := s.TokenManager.Issue(claims.UserID, claims.Subject)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenRes{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
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

	result, total, err := s.sysUserRepo.List(ctx, q)
	if err != nil {
		return res, err
	}

	items := lo.Map(result, func(item domain.SysUser, _ int) UserItem {
		return s.toUserItem(item)
	})

	return domain.NewPageResult(items, total, in.Page, in.PageSize), nil
}

func (s *Service) Create(ctx context.Context, in CreateUserReq) (res CreateUserRes, err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("creating user")

	if existing, err := s.sysUserRepo.FindByAccount(ctx, in.Account); err == nil && existing != nil {
		return res, domain.ErrAccountExists
	}

	u := &domain.SysUser{
		Account:      in.Account,
		PasswordHash: in.Password,
		NickName:     in.NickName,
		Email:        in.Email,
		Phone:        in.Phone,
		Enabled:      true,
	}
	if err := u.PwdHash(); err != nil {
		return res, fmt.Errorf("hash password for %s: %w", in.Account, err)
	}
	if len(u.NickName) == 0 {
		u.NickName = in.Account
	}

	if err := s.sysUserRepo.Create(ctx, u); err != nil {
		return res, err
	}

	if len(in.RoleIDs) > 0 {
		if err := s.sysUserRepo.CreateWithRoles(ctx, u, in.RoleIDs); err != nil {
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
	user, err := s.sysUserRepo.FindByIDWithRoles(ctx, currUser.UserID)
	if err != nil {
		return res, err
	}
	return s.toUserItem(*user), nil
}

func (s *Service) Update(ctx context.Context, in UpdateUserReq) (res UpdateUserRes, err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("update user")

	current, err := s.sysUserRepo.FindByID(ctx, in.UserID)
	if err != nil {
		return res, err
	}

	if in.NickName != nil {
		current.NickName = *in.NickName
	}
	if in.Password != nil {
		current.PasswordHash = *in.Password
		if err := current.PwdHash(); err != nil {
			return res, fmt.Errorf("hash password for %s: %w", current.Account, err)
		}
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
		if err := s.sysUserRepo.UpdateWithRoles(ctx, current, in.RoleIDs); err != nil {
			return res, err
		}
		if err := s.policy.SyncUserRolesByIDs(ctx, current.Account, in.RoleIDs); err != nil {
			return res, fmt.Errorf("sync user roles (account=%s): %w", current.Account, err)
		}
	} else {
		if err := s.sysUserRepo.Update(ctx, current); err != nil {
			return res, err
		}
	}
	return res, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	logger := s.Log(ctx)
	logger.Debug().Int64("id", id).Msg("deleting user")

	u, err := s.sysUserRepo.FindByID(ctx, id)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	if errors.Is(err, domain.ErrNotFound) {
		return nil
	}

	if s.IsAdmin(u.Account) {
		return domain.ErrCannotDeleteAdmin
	}

	if err := s.sysUserRepo.Delete(ctx, id); err != nil {
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

func (s *Service) toUserItem(u domain.SysUser) UserItem {
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
