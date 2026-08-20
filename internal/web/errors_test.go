package web

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := DecodeError(tc.err); got != tc.want {
				t.Fatalf("DecodeError() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestDecodeErrorMappings(t *testing.T) {
	t.Parallel()

	for _, mappings := range allErrorMappings {
		for _, mapping := range mappings {
			mapping := mapping
			t.Run(mapping.descriptor.Message, func(t *testing.T) {
				if got := DecodeError(mapping.target); got != mapping.descriptor {
					t.Fatalf("DecodeError(%q) = %#v, want %#v", mapping.target, got, mapping.descriptor)
				}

				wrapped := fmt.Errorf("wrapped: %w", mapping.target)
				if got := DecodeError(wrapped); got != mapping.descriptor {
					t.Fatalf("DecodeError(wrapped %q) = %#v, want %#v", mapping.target, got, mapping.descriptor)
				}
			})
		}
	}
}
