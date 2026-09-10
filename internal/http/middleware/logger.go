package middleware

import (
	"gin-layout/internal/infra/log"
	"net/http"
	"time"

	"gin-layout/internal/shared/reqctx"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			path += "?" + raw
		}

		c.Next()

		latency := time.Since(start)
		var evt *zerolog.Event
		switch {
		case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
			evt = log.Warn()
		case c.Writer.Status() >= http.StatusInternalServerError:
			evt = log.Error()
		default:
			evt = log.Info()
		}

		logEvent := evt.Ctx(c.Request.Context()).
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", path).
			Str(reqctx.RequestIDKey, c.GetString(reqctx.RequestIDKey)).
			Str("ip", c.ClientIP())

		if currentUser, ok := reqctx.CurrentUserFromContext(c.Request.Context()); ok {
			logEvent = logEvent.Int64(reqctx.UserIDKey, currentUser.UserID).Str(reqctx.UserKey, currentUser.Account)
		}

		if len(c.Errors) > 0 {
			logEvent = logEvent.Str("gin_errors", c.Errors.String())
		}

		logEvent.
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent()).
			Int("body_size", c.Writer.Size()).
			Msg("Request")
	}
}
