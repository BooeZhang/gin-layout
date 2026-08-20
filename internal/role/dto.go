package role

import (
	"time"

	"gin-layout/internal/page"
)

type ListRoleReq struct {
	page.Request
	Name    *string `json:"name" form:"name"`
	Code    *string `json:"code" form:"code"`
	Enabled *bool   `json:"enabled" form:"enabled"`
}

type CreateRoleReq struct {
	Name          string  `json:"name" binding:"required" desc:"角色名"`
	Code          string  `json:"code" binding:"required" desc:"角色编码"`
	Description   string  `json:"description" desc:"描述"`
	Sort          int     `json:"sort"`
	Enabled       bool    `json:"enabled"`
	PermissionIDs []int64 `json:"permissionIDs"`
}

type CreateRoleRes struct {
	ID int64 `json:"id"`
}

type UpdateRoleReq struct {
	RoleID        int64   `json:"-"`
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	Sort          *int    `json:"sort"`
	Enabled       *bool   `json:"enabled"`
	PermissionIDs []int64 `json:"permissionIDs"`
}

type UpdateRoleRes struct{}

type AssignPermissionsReq struct {
	PermissionIDs []int64 `json:"permission_ids" binding:"required"`
}
type AssignPermissionsRes struct{}

type RoleItem struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	Description   string    `json:"description,omitzero"`
	Sort          int       `json:"sort,omitzero"`
	CreatedAt     time.Time `json:"createdAt"`
	Enabled       bool      `json:"enabled,omitzero"`
	PermissionIDs []int64   `json:"permissionIDs,omitzero"`
}

type UserAddReq struct {
	RoleID  int64   `json:"-"`
	UserIDs []int64 `json:"userIDs" binding:"required"`
}
type UserAddRes struct{}

type UserRemoveReq struct {
	RoleID  int64   `json:"-"`
	UserIDs []int64 `json:"userIDs" binding:"required"`
}
type UserRemoveRes struct{}
