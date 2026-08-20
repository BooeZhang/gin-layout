package sysuser

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"

	"gin-layout/internal/common"
	"gin-layout/internal/domain"
	"gin-layout/internal/page"
	"gin-layout/internal/policy"
	"gin-layout/internal/reqctx"
	"gin-layout/internal/token"
)

type Service struct {
	common.BaseService

	tokens      token.Manager
	sysUserRepo Repository
	policy      policy.Manager
	roles       RoleFinder
	menus       MenuFinder
}

type Deps struct {
	Base   common.BaseService
	Tokens token.Manager
	Users  Repository
	Policy policy.Manager
	Roles  RoleFinder
	Menus  MenuFinder
}

func NewService(deps Deps) *Service {
	return &Service{
		BaseService: deps.Base,
		tokens:      deps.Tokens,
		sysUserRepo: deps.Users,
		policy:      deps.Policy,
		roles:       deps.Roles,
		menus:       deps.Menus,
	}
}

func (s *Service) Login(ctx context.Context, req LoginReq) (*LoginRes, error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", req).Msg("login attempt")

	u, found, err := s.sysUserRepo.FindByAccount(ctx, req.Account)
	if err != nil {
		return nil, fmt.Errorf("find login user: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("login: %w", ErrInvalidCredentials)
	}

	if !u.Enabled {
		return nil, ErrUserDisabled
	}

	if !u.ComparePassword(req.Password) {
		return nil, ErrInvalidCredentials
	}

	tokens, err := s.tokens.Issue(u.ID, u.Account)
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
	if _, err := s.tokens.RevokeCurrent(ctx); err != nil {
		return nil, err
	}
	return &LogoutRes{}, nil
}

func (s *Service) RefreshToken(ctx context.Context, req RefreshTokenReq) (*RefreshTokenRes, error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", req).Msg("refresh token")

	claims, err := s.tokens.Parse(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if claims.Type != token.TypeRefresh {
		return nil, token.ErrTokenInvalid
	}

	revoked, err := s.tokens.IsRevoked(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if revoked {
		return nil, token.ErrTokenRevoked
	}

	if err := s.tokens.Revoke(ctx, req.RefreshToken, claims.UserID, claims.ExpiresAt); err != nil {
		return nil, err
	}

	tokens, err := s.tokens.Issue(claims.UserID, claims.Subject)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenRes{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *Service) List(ctx context.Context, in ListUserReq) (res page.Result[UserItem], err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("user list")

	q := userListQuery{
		Request:  in.Request,
		Account:  in.Account,
		NickName: in.NickName,
		Email:    in.Email,
		Phone:    in.Phone,
		Enabled:  in.Enabled,
	}

	result, total, err := s.sysUserRepo.List(ctx, q)
	if err != nil {
		return res, err
	}

	items := lo.Map(result, func(item domain.SysUser, _ int) UserItem {
		return s.toUserItem(item)
	})

	return page.NewResult(items, total, q.Request.Page, q.Request.PageSize), nil
}

func (s *Service) Create(ctx context.Context, in CreateUserReq) (res CreateUserRes, err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("creating user")

	existing, found, err := s.sysUserRepo.FindByAccount(ctx, in.Account)
	if err != nil {
		return res, err
	}
	if found && existing != nil {
		return res, ErrAccountExists
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
	currUser, ok := reqctx.CurrentUserFromContext(ctx)
	if !ok {
		return res, token.ErrUnauthenticated
	}
	user, found, err := s.sysUserRepo.FindByIDWithRoles(ctx, currUser.UserID)
	if err != nil {
		return res, err
	}
	if !found {
		return res, ErrUserNotFound
	}
	return s.toUserItem(*user), nil
}

func (s *Service) Update(ctx context.Context, in UpdateUserReq) (res UpdateUserRes, err error) {
	logger := s.Log(ctx)
	logger.Debug().Any("input", in).Msg("update user")

	current, found, err := s.sysUserRepo.FindByID(ctx, in.UserID)
	if err != nil {
		return res, err
	}
	if !found {
		return res, ErrUserNotFound
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

	u, found, err := s.sysUserRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	if s.IsAdmin(u.Account) {
		return ErrCannotDeleteAdmin
	}

	if err := s.sysUserRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("DeleteUser (id=%d): %w", id, err)
	}
	return nil
}

func (s *Service) GetCurrentUserMenus(ctx context.Context) ([]domain.MenuItem, error) {
	currUser, ok := reqctx.CurrentUserFromContext(ctx)
	if !ok {
		return nil, token.ErrUnauthenticated
	}

	if s.IsAdmin(currUser.Account) {
		rows, err := s.menus.ListAll(ctx)
		if err != nil {
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
