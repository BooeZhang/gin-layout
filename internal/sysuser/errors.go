package sysuser

import "errors"

var (
	ErrInvalidAccountFormat = errors.New("用户名格式错")
	ErrInvalidUserID        = errors.New("无效的用户 ID")
	ErrWeakPassword         = errors.New("密码强度不足")
	ErrAccountExists        = errors.New("账号已存在")
	ErrUserNotFound         = errors.New("用户不存在")
	ErrUserDisabled         = errors.New("用户已禁用")
	ErrCannotDeleteAdmin    = errors.New("不能删除超级管理员")
	ErrInvalidCredentials   = errors.New("用户名或密码错误")
	ErrPasswordIdentical    = errors.New("新密码不能与旧密码相同")
)
