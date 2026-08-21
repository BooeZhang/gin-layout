package sysuser

import "gin-layout/internal/apperror"

var (
	ErrInvalidAccountFormat = apperror.New(apperror.InvalidInput, 20010, "用户名格式错")
	ErrInvalidUserID        = apperror.New(apperror.InvalidInput, 20011, "无效的用户 ID")
	ErrWeakPassword         = apperror.New(apperror.InvalidInput, 20012, "密码强度不足")
	ErrAccountExists        = apperror.New(apperror.Conflict, 20020, "账号已存在")
	ErrUserNotFound         = apperror.New(apperror.NotFound, 20000, "用户不存在")
	ErrUserDisabled         = apperror.New(apperror.BusinessResult, 20040, "用户已禁用")
	ErrCannotDeleteAdmin    = apperror.New(apperror.BusinessResult, 20050, "不能删除超级管理员")
	ErrInvalidCredentials   = apperror.New(apperror.BusinessResult, 20051, "用户名或密码错误")
	ErrPasswordIdentical    = apperror.New(apperror.BusinessResult, 20052, "新密码不能与旧密码相同")
)
