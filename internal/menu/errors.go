package menu

import "gin-layout/internal/apperror"

var (
	ErrInvalidMenuID = apperror.New(apperror.InvalidInput, 40010, "无效ID")
	ErrMenuExists    = apperror.New(apperror.Conflict, 40020, "菜单已存在")
	ErrMenuNotFound  = apperror.New(apperror.NotFound, 40000, "菜单不存在")
)
