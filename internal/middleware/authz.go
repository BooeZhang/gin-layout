package middleware

import (
	"github.com/BooeZhang/gin-layout/pkg/jwtx"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewAuthorizer(e *casbin.Enforcer) gin.HandlerFunc {
	a := &BasicAuthorizer{enforcer: e}
	return func(c *gin.Context) {
		if !a.CheckPermission(c) {
			a.RequirePermission(c)
		}
		c.Next()
	}
}

type BasicAuthorizer struct {
	enforcer *casbin.Enforcer
}

func (a *BasicAuthorizer) GetUserId(c *gin.Context) string {
	claims, _ := c.Get("userClaims")
	userClaims, ok := claims.(jwtx.UserClaims)
	if ok {
		return strconv.FormatInt(int64(userClaims.UserId), 10)
	}
	return ""
}

func (a *BasicAuthorizer) CheckPermission(c *gin.Context) bool {
	_ = a.enforcer.LoadPolicy()
	userId := a.GetUserId(c)
	method := c.Request.Method
	path := c.Request.RequestURI

	allowed, err := a.enforcer.Enforce(userId, path, method)
	if err != nil {
		log.L(c).Errorf("校验权限失败: %s", err)
		return false
	}

	return allowed
}

func (a *BasicAuthorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(http.StatusForbidden)
}
