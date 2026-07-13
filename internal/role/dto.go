package role

import "time"

type ListRoleReq struct {
	Page    int     `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Name    *string `json:"name" form:"name"`
	Code    *string `json:"code" form:"code"`
	Enabled *bool   `json:"enabled" form:"enabled"`
}

type CreateRoleReq struct {
	Name          string  `json:"name" binding:"required"`
	Code          string  `json:"code" binding:"required"`
	Description   string  `json:"description"`
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
	Description   string    `json:"description,omitempty"`
	Sort          int       `json:"sort,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	Enabled       bool      `json:"enabled,omitempty"`
	PermissionIDs []int64   `json:"permissionIDs,omitempty"`
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
