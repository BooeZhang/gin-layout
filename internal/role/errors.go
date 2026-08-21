package role

import "gin-layout/internal/apperror"

var (
	ErrInvalidRoleID         = apperror.New(apperror.InvalidInput, 30010, "无效ID")
	ErrRoleExists            = apperror.New(apperror.Conflict, 30020, "角色已存在")
	ErrRoleNotFound          = apperror.New(apperror.NotFound, 30000, "角色不存在")
	ErrPermissionNotFound    = apperror.New(apperror.NotFound, 30100, "权限不存在")
	ErrRoleDisabled          = apperror.New(apperror.BusinessResult, 30040, "角色已禁用")
	ErrCannotDeleteAdminRole = apperror.New(apperror.BusinessResult, 30050, "不允许删除管理员角色")
)
