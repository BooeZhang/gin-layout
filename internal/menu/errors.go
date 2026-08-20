package menu

import "errors"

var (
	ErrInvalidMenuID = errors.New("无效ID")
	ErrMenuExists    = errors.New("菜单已存在")
	ErrMenuNotFound  = errors.New("菜单不存在")
)
