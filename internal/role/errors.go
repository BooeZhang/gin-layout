package role

import "errors"

var (
	ErrInvalidRoleID         = errors.New("无效ID")
	ErrRoleExists            = errors.New("角色已存在")
	ErrRoleNotFound          = errors.New("角色不存在")
	ErrPermissionNotFound    = errors.New("权限不存在")
	ErrRoleDisabled          = errors.New("角色已禁用")
	ErrCannotDeleteAdminRole = errors.New("不允许删除管理员角色")
)
