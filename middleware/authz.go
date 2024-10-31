package middleware

import (
	"net/http"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"gin-layout/pkg/jwtx"
)

type BasicAuthorizer struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizer(e *casbin.Enforcer) gin.HandlerFunc {
	a := &BasicAuthorizer{enforcer: e}
	return func(ctx *gin.Context) {
		if !a.CheckPermission(ctx) {
			a.RequirePermission(ctx)
		}
		ctx.Next()
	}
}

func (a *BasicAuthorizer) CheckPermission(ctx *gin.Context) bool {
	_ = a.enforcer.LoadPolicy()
	userId := a.GetUserId(ctx)
	method := ctx.Request.Method
	path := ctx.Request.URL.Path

	allowed, err := a.enforcer.Enforce(userId, path, method)
	if err != nil {
		log.Error().Err(err).Msg("校验权限失败")
		return false
	}

	return allowed
}

func (a *BasicAuthorizer) GetUserId(ctx *gin.Context) string {
	userId := jwtx.GetUserID(ctx)
	return strconv.FormatInt(int64(userId), 10)
}

func (a *BasicAuthorizer) RequirePermission(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "msg": "没有权限"})
}
