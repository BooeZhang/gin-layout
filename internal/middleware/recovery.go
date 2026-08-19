package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"gin-layout/internal/infra"
	"gin-layout/internal/web"
)

func Recovery(logger *infra.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				l := infra.LogFromContext(c.Request.Context(), logger)

				err := fmt.Errorf("panic recovered: %v", r)
				l.Error().
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
