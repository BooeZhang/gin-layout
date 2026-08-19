package menu

import (
	"net/http"

	"gin-layout/internal/domain"
)

var (
	ErrInvalidMenuID = domain.NewDomainError(40010, http.StatusUnprocessableEntity, "无效ID")
	ErrMenuExists    = domain.NewDomainError(40020, http.StatusConflict, "菜单已存在")
)
