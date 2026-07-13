package domain

import "net/http"

// ── Common (1xxxx) ──

var (
	ErrInvalidArgument = NewDomainError(10010, http.StatusUnprocessableEntity, "invalid argument")
	ErrNotFound        = NewDomainError(10000, http.StatusNotFound, "resource not found")
)

// ── User/Account (2xxxx) ──

var (
	ErrInvalidAccountFormat = NewDomainError(20010, http.StatusUnprocessableEntity, "用户名格式错")
	ErrInvalidUserID        = NewDomainError(20011, http.StatusUnprocessableEntity, "无效的用户 ID")
	ErrWeakPassword         = NewDomainError(20012, http.StatusUnprocessableEntity, "密码强度不足")
	ErrAccountExists        = NewDomainError(20020, http.StatusConflict, "账号已存在")
	ErrUserNotExist         = NewDomainError(20000, http.StatusNotFound, "用户不存在")
	ErrUserDisabled         = NewDomainError(20040, http.StatusOK, "用户已禁用")
	ErrCannotDeleteAdmin    = NewDomainError(20050, http.StatusOK, "不能删除超级管理员")
	ErrAccountOrPassword    = NewDomainError(20051, http.StatusOK, "用户名或密码错误")
	ErrPasswordIdentical    = NewDomainError(20052, http.StatusOK, "新密码不能与旧密码相同")
)

// ── Role/Permission (3xxxx) ──

var (
	ErrInvalidRoleID         = NewDomainError(30010, http.StatusUnprocessableEntity, "无效ID")
	ErrRoleExists            = NewDomainError(30020, http.StatusConflict, "角色已存在")
	ErrPermNotExist          = NewDomainError(30100, http.StatusNotFound, "权限不存在")
	ErrRoleDisabled          = NewDomainError(30040, http.StatusOK, "角色已禁用")
	ErrCannotDeleteAdminRole = NewDomainError(30050, http.StatusOK, "不允许删除管理员角色")
	ErrNotPermission         = NewDomainError(30130, http.StatusForbidden, "没有权限")
)

// ── Menu/Resource (4xxxx) ──

var (
	ErrInvalidMenuID = NewDomainError(40010, http.StatusUnprocessableEntity, "无效ID")
	ErrMenuExists    = NewDomainError(40020, http.StatusConflict, "菜单已存在")
)

// ── Auth/Token (5xxxx) ──

var (
	ErrInvalidAccessToken = NewDomainError(50010, http.StatusUnauthorized, "无效访问令牌")
	ErrNotLogin           = NewDomainError(50011, http.StatusUnauthorized, "未登录或非法访问")
	ErrTokenInvalid       = NewDomainError(50012, http.StatusUnauthorized, "token 无效")
	ErrTokenExpired       = NewDomainError(50060, http.StatusUnauthorized, "token 已过期")
	ErrTokenNotActive     = NewDomainError(50070, http.StatusUnauthorized, "token 不是活跃状态")
	ErrTokenRevoked       = NewDomainError(50061, http.StatusUnauthorized, "token 已失效")
)
