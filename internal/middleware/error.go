package middleware

import (
	"github.com/gin-gonic/gin"

	"gin-layout/internal/infra"
	"gin-layout/internal/web"
)

func ErrorHandler(logger *infra.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		l := infra.LogFromContext(c.Request.Context(), logger)

		descriptor := web.DecodeError(err)
		l.Error().
			Err(err).
			Int("biz_code", descriptor.Code).
			Int("http_status", descriptor.HTTPStatus).
			Str("biz_message", descriptor.Message).
			Msg("request failed")

		web.Error(c, err)
		c.Abort()
	}
}
