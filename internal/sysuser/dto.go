package sysuser

import (
	"time"

	"gin-layout/internal/page"
)

type LoginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type (
	LogoutReq struct{}
	LogoutRes struct{}
)

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type CreateUserReq struct {
	Account  string  `json:"account" binding:"required,account"`
	Password string  `json:"password" binding:"required,password"`
	NickName string  `json:"nick_name"`
	Email    string  `json:"email" binding:"omitzero,email"`
	Phone    string  `json:"phone"`
	RoleIDs  []int64 `json:"roleIds"`
}

type CreateUserRes struct {
	ID int64 `json:"id"`
}

type ListUserReq struct {
	page.Request
	Account  *string `json:"account" form:"account"`
	NickName *string `json:"nick_name" form:"nick_name"`
	Email    *string `json:"email" form:"email"`
	Phone    *string `json:"phone" form:"phone"`
	Enabled  *bool   `json:"enabled" form:"enabled"`
}

type UserItem struct {
	ID          int64      `json:"id"`
	Account     string     `json:"account"`
	NickName    string     `json:"nickName"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Avatar      string     `json:"avatar"`
	Enabled     bool       `json:"enabled"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Roles       []RoleItem `json:"roles,omitzero"`
}

type RoleItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type UpdateUserReq struct {
	UserID   int64   `json:"_"`
	NickName *string `json:"nick_name"`
	Password *string `json:"password" binding:"omitzero,password"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Avatar   *string `json:"avatar"`
	Enabled  *bool   `json:"enabled"`
	RoleIDs  []int64 `json:"roleIds,omitzero"`
}

type UpdateUserRes struct{}

type ReplaceRolesReq struct {
	RoleIDs []int64 `json:"role_ids" binding:"required"`
}
