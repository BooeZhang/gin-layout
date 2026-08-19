package role

import (
	"net/http"

	"gin-layout/internal/domain"
)

var (
	ErrInvalidRoleID         = domain.NewDomainError(30010, http.StatusUnprocessableEntity, "无效ID")
	ErrRoleExists            = domain.NewDomainError(30020, http.StatusConflict, "角色已存在")
	ErrPermNotExist          = domain.NewDomainError(30100, http.StatusNotFound, "权限不存在")
	ErrRoleDisabled          = domain.NewDomainError(30040, http.StatusOK, "角色已禁用")
	ErrCannotDeleteAdminRole = domain.NewDomainError(30050, http.StatusOK, "不允许删除管理员角色")
)
