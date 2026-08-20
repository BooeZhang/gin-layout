package web

import (
	"net/http"

	"gin-layout/internal/menu"
)

var menuErrorMappings = []errorMapping{
	newErrorMapping(menu.ErrInvalidMenuID, http.StatusUnprocessableEntity, 40010, "无效ID"),
	newErrorMapping(menu.ErrMenuExists, http.StatusConflict, 40020, "菜单已存在"),
	newErrorMapping(menu.ErrMenuNotFound, http.StatusNotFound, 40000, "菜单不存在"),
}
