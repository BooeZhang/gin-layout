package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"gin-layout/internal/common"
	"gin-layout/internal/domain"
)

func RBAC(enforcer common.PolicyManager, permissions common.PermissionResolver, logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if enforcer == nil {
			c.Error(domain.ErrNotLogin)
			return
		}

		userAccount, ok := domain.CurrentUserFromContext(c.Request.Context())
		if !ok {
			c.Error(domain.ErrNotLogin)
			return
		}
		requestID, _ := domain.RequestIDFromContext(c.Request.Context())

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
			c.Error(err)
			return
		}

		if !pass {
			c.Error(domain.ErrNotPermission)
			return
		}

		c.Next()
	}
}
