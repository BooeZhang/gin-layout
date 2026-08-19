package sysuser

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"gin-layout/config"
	"gin-layout/internal/common"
	"gin-layout/internal/domain"
	"gin-layout/internal/token"
)

type fakeUserRepo struct {
	user  *domain.SysUser
	saved *domain.SysUser
}

func (f *fakeUserRepo) Create(ctx context.Context, entity *domain.SysUser) error { return nil }

func (f *fakeUserRepo) Update(ctx context.Context, entity *domain.SysUser) error {
	f.saved = entity
	return nil
}

func (f *fakeUserRepo) Delete(ctx context.Context, id int64) error { return nil }

func (f *fakeUserRepo) FindByID(ctx context.Context, id int64) (*domain.SysUser, error) {
	if f.user == nil {
		return nil, domain.ErrNotFound
	}
	return f.user, nil
}

func (f *fakeUserRepo) List(ctx context.Context, q userListQuery) ([]domain.SysUser, int64, error) {
	return nil, 0, nil
}

func (f *fakeUserRepo) FindByAccount(ctx context.Context, account string) (*domain.SysUser, error) {
	return nil, domain.ErrNotFound
}

func (f *fakeUserRepo) FindByIDWithRoles(ctx context.Context, id int64) (*domain.SysUser, error) {
	return f.FindByID(ctx, id)
}

func (f *fakeUserRepo) UpdateLastLogin(ctx context.Context, userID int64, lastLoginAt time.Time) error {
	return nil
}

func (f *fakeUserRepo) CreateWithRoles(ctx context.Context, u *domain.SysUser, roleIDs []int64) error {
	return nil
}

func (f *fakeUserRepo) UpdateWithRoles(ctx context.Context, u *domain.SysUser, roleIDs []int64) error {
	return nil
}

func (f *fakeUserRepo) ReplaceUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	return nil
}

func (f *fakeUserRepo) FindByIDs(ctx context.Context, ids []int64) ([]domain.SysUser, error) {
	return nil, nil
}

type fakePolicyManager struct{}

func (fakePolicyManager) SyncUserRoles(ctx context.Context, userAccount string, roleCodes []string) error {
	return nil
}

func (fakePolicyManager) SyncUserRolesByIDs(ctx context.Context, userAccount string, roleIDs []int64) error {
	return nil
}

func (fakePolicyManager) SyncRolePermissions(ctx context.Context, roleCode string, permissions [][]string) error {
	return nil
}

func (fakePolicyManager) AddRoleToUser(ctx context.Context, userAccount string, roleCode string) error {
	return nil
}

func (fakePolicyManager) DeleteRole(ctx context.Context, roleCode string) error {
	return nil
}

func (fakePolicyManager) Enforce(subject, object, action string) (bool, error) {
	return true, nil
}

type fakeTokenManager struct {
	claims *token.Claims
}

func (f fakeTokenManager) Issue(userID int64, subject string) (token.Pair, error) {
	return token.Pair{}, nil
}

func (f fakeTokenManager) Parse(raw string) (*token.Claims, error) {
	return f.claims, nil
}

func (f fakeTokenManager) IsRevoked(ctx context.Context, raw string) (bool, error) {
	return false, nil
}

func (f fakeTokenManager) Revoke(ctx context.Context, raw string, userID int64, expiresAt time.Time) error {
	return nil
}

func (f fakeTokenManager) RevokeCurrent(ctx context.Context) (bool, error) {
	return false, nil
}

type fakeRoleFinder struct{}

func (fakeRoleFinder) ListEnabledRoleIDsForUser(ctx context.Context, userID int64) ([]int64, error) {
	return nil, nil
}

type fakeMenuFinder struct{}

func (fakeMenuFinder) ListEnabledByRoleIDs(ctx context.Context, roleIDs []int64) ([]domain.Menu, error) {
	return nil, nil
}

func (fakeMenuFinder) ListAll(ctx context.Context) ([]domain.Menu, error) {
	return nil, nil
}

func (fakeMenuFinder) ToMenuTree(rows []domain.Menu) []domain.MenuItem {
	return nil
}

func newTestService(repo Repository, tokens ...token.Manager) *Service {
	var tokenManager token.Manager = fakeTokenManager{
		claims: &token.Claims{Type: token.TypeRefresh},
	}
	if len(tokens) > 0 {
		tokenManager = tokens[0]
	}

	logger := zerolog.Nop()
	return &Service{
		BaseService: common.NewBaseService(&config.Config{}, &logger),
		tokens:      tokenManager,
		sysUserRepo: repo,
		policy:      fakePolicyManager{},
		roles:       fakeRoleFinder{},
		menus:       fakeMenuFinder{},
	}
}

func TestService_Update_PersistsHashedPassword(t *testing.T) {
	repo := &fakeUserRepo{user: &domain.SysUser{
		ID:           1,
		Account:      "alice",
		PasswordHash: "old-hash",
	}}
	svc := newTestService(repo)

	req := UpdateUserReq{UserID: 1, Password: strPtr("new-secret")}
	if _, err := svc.Update(t.Context(), req); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if repo.saved == nil {
		t.Fatal("Update() did not persist user")
	}
	if repo.saved.PasswordHash == "old-hash" {
		t.Fatal("Update() persisted the old password hash")
	}
	if !repo.saved.ComparePassword("new-secret") {
		t.Fatal("Update() persisted password does not match new password")
	}
}

func TestService_RefreshToken_RejectsAccessToken(t *testing.T) {
	repo := &fakeUserRepo{}
	tokens := fakeTokenManager{claims: &token.Claims{Type: token.TypeAccess}}
	svc := newTestService(repo, tokens)

	_, err := svc.RefreshToken(t.Context(), RefreshTokenReq{RefreshToken: "access-token"})
	if !errors.Is(err, token.ErrTokenInvalid) {
		t.Fatalf("RefreshToken() error = %v, want %v", err, token.ErrTokenInvalid)
	}
}

func strPtr(s string) *string {
	return &s
}
