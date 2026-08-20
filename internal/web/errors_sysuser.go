package web

import (
	"net/http"

	"gin-layout/internal/sysuser"
)

var sysuserErrorMappings = []errorMapping{
	newErrorMapping(sysuser.ErrInvalidAccountFormat, http.StatusUnprocessableEntity, 20010, "用户名格式错"),
	newErrorMapping(sysuser.ErrInvalidUserID, http.StatusUnprocessableEntity, 20011, "无效的用户 ID"),
	newErrorMapping(sysuser.ErrWeakPassword, http.StatusUnprocessableEntity, 20012, "密码强度不足"),
	newErrorMapping(sysuser.ErrAccountExists, http.StatusConflict, 20020, "账号已存在"),
	newErrorMapping(sysuser.ErrUserNotFound, http.StatusNotFound, 20000, "用户不存在"),
	newErrorMapping(sysuser.ErrUserDisabled, http.StatusOK, 20040, "用户已禁用"),
	newErrorMapping(sysuser.ErrCannotDeleteAdmin, http.StatusOK, 20050, "不能删除超级管理员"),
	newErrorMapping(sysuser.ErrInvalidCredentials, http.StatusOK, 20051, "用户名或密码错误"),
	newErrorMapping(sysuser.ErrPasswordIdentical, http.StatusOK, 20052, "新密码不能与旧密码相同"),
}
