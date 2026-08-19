package common

import (
	"context"
	"time"

	"gin-layout/internal/domain"
)

type TokenIssuer interface {
	Issue(userID int64, subject string) (domain.TokenPair, error)
	Parse(raw string) (*domain.TokenClaims, error)
}

type TokenBlacklistRepository interface {
	Exists(ctx context.Context, tokenHash string) (bool, error)
	Add(ctx context.Context, tokenHash string, userID int64, expiresAt time.Time) error
}

type TokenManager interface {
	TokenIssuer
	IsRevoked(ctx context.Context, raw string) (bool, error)
	Revoke(ctx context.Context, raw string, userID int64, expiresAt time.Time) error
	RevokeCurrent(ctx context.Context) (bool, error)
}

type PolicyManager interface {
	SyncUserRoles(ctx context.Context, userAccount string, roleCodes []string) error
	SyncUserRolesByIDs(ctx context.Context, userAccount string, roleIDs []int64) error
	SyncRolePermissions(ctx context.Context, roleCode string, permissions [][]string) error
	AddRoleToUser(ctx context.Context, userAccount string, roleCode string) error
	DeleteRole(ctx context.Context, roleCode string) error
	Enforce(subject, object, action string) (bool, error)
}

type PermissionResolver interface {
	ResolvePermissionCode(url, method string) (string, bool)
}

type PageReq struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}
