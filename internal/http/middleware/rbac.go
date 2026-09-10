package middleware

import (
	"gin-layout/internal/infra/log"

	"gin-layout/internal/shared/policy"
	"gin-layout/internal/shared/reqctx"
	"gin-layout/internal/shared/token"

	"github.com/gin-gonic/gin"
)

func RBAC(enforcer policy.Manager, permissions policy.PermissionResolver) gin.HandlerFunc {
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
			log.Ctx(c.Request.Context()).Error().Err(err).
				Str("userAccount", userAccount.Account).
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
