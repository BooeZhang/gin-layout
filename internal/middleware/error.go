package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"

	"gin-layout/internal/domain"
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

		var bizErr *domain.DomainError
		if errors.As(err, &bizErr) {
			l.Error().
				Err(err).
				Int("biz_code", bizErr.Code).
				Int("http_status", bizErr.HTTPStatus).
				Str("biz_message", bizErr.Message).
				Msg("request failed")
		} else {
			l.Error().
				Err(err).
				Msg("request failed with unknown error")
		}

		web.Error(c, err)
		c.Abort()
	}
}
