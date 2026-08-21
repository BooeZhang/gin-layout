package web

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"gin-layout/internal/apperror"
	"gin-layout/internal/menu"
	"gin-layout/internal/policy"
	"gin-layout/internal/role"
	"gin-layout/internal/sysuser"
	"gin-layout/internal/token"
)

func TestDecodeError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		err  error
		want ErrorDescriptor
	}{
		{
			name: "sysuser not found",
			err:  sysuser.ErrUserNotFound,
			want: ErrorDescriptor{HTTPStatus: http.StatusNotFound, Code: 20000, Message: "用户不存在"},
		},
		{
			name: "wrapped unauthenticated",
			err:  fmt.Errorf("authenticate request: %w", token.ErrUnauthenticated),
			want: ErrorDescriptor{HTTPStatus: http.StatusUnauthorized, Code: 50011, Message: "未登录或非法访问"},
		},
		{
			name: "role permission not found",
			err:  role.ErrPermissionNotFound,
			want: ErrorDescriptor{HTTPStatus: http.StatusNotFound, Code: 30100, Message: "权限不存在"},
		},
		{
			name: "permission denied",
			err:  policy.ErrPermissionDenied,
			want: ErrorDescriptor{HTTPStatus: http.StatusForbidden, Code: 30130, Message: "没有权限"},
		},
		{
			name: "unknown error",
			err:  errors.New("database unavailable"),
			want: ErrorDescriptor{HTTPStatus: http.StatusInternalServerError, Code: 50001, Message: "internal server error"},
		},
		{
			name: "structured not found",
			err:  apperror.New(apperror.NotFound, 20999, "资源不存在"),
			want: ErrorDescriptor{HTTPStatus: http.StatusNotFound, Code: 20999, Message: "资源不存在"},
		},
		{
			name: "structured business error",
			err:  apperror.New(apperror.BusinessResult, 20998, "业务操作失败"),
			want: ErrorDescriptor{HTTPStatus: http.StatusOK, Code: 20998, Message: "业务操作失败"},
		},
		{
			name: "unknown kind",
			err:  apperror.New(apperror.Unknown, 50002, "未知错误"),
			want: ErrorDescriptor{HTTPStatus: http.StatusInternalServerError, Code: 50002, Message: "未知错误"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := DecodeError(tc.err); got != tc.want {
				t.Fatalf("DecodeError() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestDecodeErrorDomainErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		err  error
		want ErrorDescriptor
	}{
		{name: "invalid menu id", err: menu.ErrInvalidMenuID, want: ErrorDescriptor{HTTPStatus: http.StatusUnprocessableEntity, Code: 40010, Message: "无效ID"}},
		{name: "menu exists", err: menu.ErrMenuExists, want: ErrorDescriptor{HTTPStatus: http.StatusConflict, Code: 40020, Message: "菜单已存在"}},
		{name: "menu not found", err: menu.ErrMenuNotFound, want: ErrorDescriptor{HTTPStatus: http.StatusNotFound, Code: 40000, Message: "菜单不存在"}},
		{name: "invalid role id", err: role.ErrInvalidRoleID, want: ErrorDescriptor{HTTPStatus: http.StatusUnprocessableEntity, Code: 30010, Message: "无效ID"}},
		{name: "role exists", err: role.ErrRoleExists, want: ErrorDescriptor{HTTPStatus: http.StatusConflict, Code: 30020, Message: "角色已存在"}},
		{name: "role not found", err: role.ErrRoleNotFound, want: ErrorDescriptor{HTTPStatus: http.StatusNotFound, Code: 30000, Message: "角色不存在"}},
		{name: "permission not found", err: role.ErrPermissionNotFound, want: ErrorDescriptor{HTTPStatus: http.StatusNotFound, Code: 30100, Message: "权限不存在"}},
		{name: "role disabled", err: role.ErrRoleDisabled, want: ErrorDescriptor{HTTPStatus: http.StatusOK, Code: 30040, Message: "角色已禁用"}},
		{name: "cannot delete admin role", err: role.ErrCannotDeleteAdminRole, want: ErrorDescriptor{HTTPStatus: http.StatusOK, Code: 30050, Message: "不允许删除管理员角色"}},
		{name: "invalid account format", err: sysuser.ErrInvalidAccountFormat, want: ErrorDescriptor{HTTPStatus: http.StatusUnprocessableEntity, Code: 20010, Message: "用户名格式错"}},
		{name: "invalid user id", err: sysuser.ErrInvalidUserID, want: ErrorDescriptor{HTTPStatus: http.StatusUnprocessableEntity, Code: 20011, Message: "无效的用户 ID"}},
		{name: "weak password", err: sysuser.ErrWeakPassword, want: ErrorDescriptor{HTTPStatus: http.StatusUnprocessableEntity, Code: 20012, Message: "密码强度不足"}},
		{name: "account exists", err: sysuser.ErrAccountExists, want: ErrorDescriptor{HTTPStatus: http.StatusConflict, Code: 20020, Message: "账号已存在"}},
		{name: "user not found", err: sysuser.ErrUserNotFound, want: ErrorDescriptor{HTTPStatus: http.StatusNotFound, Code: 20000, Message: "用户不存在"}},
		{name: "user disabled", err: sysuser.ErrUserDisabled, want: ErrorDescriptor{HTTPStatus: http.StatusOK, Code: 20040, Message: "用户已禁用"}},
		{name: "cannot delete admin", err: sysuser.ErrCannotDeleteAdmin, want: ErrorDescriptor{HTTPStatus: http.StatusOK, Code: 20050, Message: "不能删除超级管理员"}},
		{name: "invalid credentials", err: sysuser.ErrInvalidCredentials, want: ErrorDescriptor{HTTPStatus: http.StatusOK, Code: 20051, Message: "用户名或密码错误"}},
		{name: "password identical", err: sysuser.ErrPasswordIdentical, want: ErrorDescriptor{HTTPStatus: http.StatusOK, Code: 20052, Message: "新密码不能与旧密码相同"}},
		{name: "invalid access token", err: token.ErrInvalidAccessToken, want: ErrorDescriptor{HTTPStatus: http.StatusUnauthorized, Code: 50010, Message: "无效访问令牌"}},
		{name: "unauthenticated", err: token.ErrUnauthenticated, want: ErrorDescriptor{HTTPStatus: http.StatusUnauthorized, Code: 50011, Message: "未登录或非法访问"}},
		{name: "token invalid", err: token.ErrTokenInvalid, want: ErrorDescriptor{HTTPStatus: http.StatusUnauthorized, Code: 50012, Message: "token 无效"}},
		{name: "token expired", err: token.ErrTokenExpired, want: ErrorDescriptor{HTTPStatus: http.StatusUnauthorized, Code: 50060, Message: "token 已过期"}},
		{name: "token revoked", err: token.ErrTokenRevoked, want: ErrorDescriptor{HTTPStatus: http.StatusUnauthorized, Code: 50061, Message: "token 已失效"}},
		{name: "token not active", err: token.ErrTokenNotActive, want: ErrorDescriptor{HTTPStatus: http.StatusUnauthorized, Code: 50070, Message: "token 不是活跃状态"}},
		{name: "permission denied", err: policy.ErrPermissionDenied, want: ErrorDescriptor{HTTPStatus: http.StatusForbidden, Code: 30130, Message: "没有权限"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := DecodeError(tc.err); got != tc.want {
				t.Fatalf("DecodeError(%q) = %#v, want %#v", tc.err, got, tc.want)
			}

			wrapped := fmt.Errorf("wrapped: %w", tc.err)
			if got := DecodeError(wrapped); got != tc.want {
				t.Fatalf("DecodeError(wrapped %q) = %#v, want %#v", tc.err, got, tc.want)
			}
		})
	}
}
