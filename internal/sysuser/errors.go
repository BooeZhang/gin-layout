package sysuser

import (
	"net/http"

	"gin-layout/internal/domain"
)

var (
	ErrInvalidAccountFormat = domain.NewDomainError(20010, http.StatusUnprocessableEntity, "用户名格式错")
	ErrInvalidUserID        = domain.NewDomainError(20011, http.StatusUnprocessableEntity, "无效的用户 ID")
	ErrWeakPassword         = domain.NewDomainError(20012, http.StatusUnprocessableEntity, "密码强度不足")
	ErrAccountExists        = domain.NewDomainError(20020, http.StatusConflict, "账号已存在")
	ErrUserNotExist         = domain.NewDomainError(20000, http.StatusNotFound, "用户不存在")
	ErrUserDisabled         = domain.NewDomainError(20040, http.StatusOK, "用户已禁用")
	ErrCannotDeleteAdmin    = domain.NewDomainError(20050, http.StatusOK, "不能删除超级管理员")
	ErrAccountOrPassword    = domain.NewDomainError(20051, http.StatusOK, "用户名或密码错误")
	ErrPasswordIdentical    = domain.NewDomainError(20052, http.StatusOK, "新密码不能与旧密码相同")
)
