package policy

import (
	"context"

	"gin-layout/internal/apperror"
)

var ErrPermissionDenied = apperror.New(apperror.Forbidden, 30130, "没有权限")

type Manager interface {
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
