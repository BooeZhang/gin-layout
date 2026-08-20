package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"gin-layout/internal/policy"
	"gin-layout/internal/reqctx"
	"gin-layout/internal/token"
)

func RBAC(enforcer policy.Manager, permissions policy.PermissionResolver, logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if enforcer == nil {
			_ = c.Error(token.ErrUnauthenticated)
			return
		}

		userAccount, ok := reqctx.CurrentUserFromContext(c.Request.Context())
		if !ok {
			_ = c.Error(token.ErrUnauthenticated)
			return
		}
		requestID, _ := reqctx.RequestIDFromContext(c.Request.Context())

		permissionCode, ok := "", false
		if permissions != nil {
			permissionCode, ok = permissions.ResolvePermissionCode(c.FullPath(), c.Request.Method)
			if !ok {
				permissionCode, ok = permissions.ResolvePermissionCode(c.Request.URL.Path, c.Request.Method)
			}
		}
		if !ok {
			permissionCode = c.Request.URL.Path
		}

		pass, err := enforcer.Enforce(userAccount.Account, permissionCode, "exec")
		if err != nil {
			logger.Error().Err(err).
				Str("userAccount", userAccount.Account).
				Str("requestID", requestID).
				Str("url", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("permissionCode", permissionCode).
				Msg("enforcer.Enforce failed")
			_ = c.Error(err)
			return
		}

		if !pass {
			_ = c.Error(policy.ErrPermissionDenied)
			return
		}

		c.Next()
	}
}
