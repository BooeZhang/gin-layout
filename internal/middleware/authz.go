package middleware

import (
	"github.com/BooeZhang/gin-layout/pkg/jwtx"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
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

func (a *BasicAuthorizer) GetRoleName(c *gin.Context) string {
	claims, _ := c.Get("userClaims")
	userClaims, ok := claims.(jwtx.UserClaims)
	if ok {
		return userClaims.Role
	}
	return ""
}

func (a *BasicAuthorizer) CheckPermission(c *gin.Context) bool {
	role := a.GetRoleName(c)
	method := c.Request.Method
	path := c.Request.RequestURI

	allowed, err := a.enforcer.Enforce(role, path, method)
	if err != nil {
		panic(err)
	}

	return allowed
}

func (a *BasicAuthorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(http.StatusForbidden)
}
