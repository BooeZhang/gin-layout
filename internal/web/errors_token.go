package web

import (
	"net/http"

	"gin-layout/internal/token"
)

var tokenErrorMappings = []errorMapping{
	newErrorMapping(token.ErrInvalidAccessToken, http.StatusUnauthorized, 50010, "无效访问令牌"),
	newErrorMapping(token.ErrUnauthenticated, http.StatusUnauthorized, 50011, "未登录或非法访问"),
	newErrorMapping(token.ErrTokenInvalid, http.StatusUnauthorized, 50012, "token 无效"),
	newErrorMapping(token.ErrTokenExpired, http.StatusUnauthorized, 50060, "token 已过期"),
	newErrorMapping(token.ErrTokenRevoked, http.StatusUnauthorized, 50061, "token 已失效"),
	newErrorMapping(token.ErrTokenNotActive, http.StatusUnauthorized, 50070, "token 不是活跃状态"),
}
