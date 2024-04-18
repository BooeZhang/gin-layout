package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/BooeZhang/gin-layout/pkg/log"
)

// UserKey defines the key in gin context which represents the owner of the secret.
const (
	UserKey = "user"
)

// Context is a middleware that injects common prefix fields to gin.Context.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(log.KeyRequestID, c.GetString(XRequestIDKey))
		c.Set(log.KeyUser, c.GetString(UserKey))
		c.Next()
	}
}
