package web

import (
	"net/http"

	"gin-layout/internal/role"
)

var roleErrorMappings = []errorMapping{
	newErrorMapping(role.ErrInvalidRoleID, http.StatusUnprocessableEntity, 30010, "无效ID"),
	newErrorMapping(role.ErrRoleExists, http.StatusConflict, 30020, "角色已存在"),
	newErrorMapping(role.ErrRoleNotFound, http.StatusNotFound, 30000, "角色不存在"),
	newErrorMapping(role.ErrPermissionNotFound, http.StatusNotFound, 30100, "权限不存在"),
	newErrorMapping(role.ErrRoleDisabled, http.StatusOK, 30040, "角色已禁用"),
	newErrorMapping(role.ErrCannotDeleteAdminRole, http.StatusOK, 30050, "不允许删除管理员角色"),
}
