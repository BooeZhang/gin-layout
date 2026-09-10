package middleware

import (
	"fmt"
	"gin-layout/internal/infra/log"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"gin-layout/internal/http/web"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("panic recovered: %v", r)
				log.Ctx(c.Request.Context()).Error().
					Err(err).
					Str("stack", string(debug.Stack())).
					Msg("panic recovered")

				web.Error(c, err)
				c.Abort()
			}
		}()

		c.Next()
	}
}
