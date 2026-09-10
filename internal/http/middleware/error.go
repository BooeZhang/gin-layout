package middleware

import (
	"gin-layout/internal/infra/log"

	"github.com/gin-gonic/gin"

	"gin-layout/internal/http/web"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		descriptor := web.DecodeError(err)
		log.Ctx(c.Request.Context()).Error().
			Err(err).
			Int("biz_code", descriptor.Code).
			Int("http_status", descriptor.HTTPStatus).
			Str("biz_message", descriptor.Message).
			Msg("request failed")

		web.Error(c, err)
		c.Abort()
	}
}
