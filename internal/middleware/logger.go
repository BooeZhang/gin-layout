package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"gin-layout/internal/domain"
	"gin-layout/internal/infra"
)

func RequestLogger(logger *infra.Logger) gin.HandlerFunc {
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
			evt = logger.Warn()
		case c.Writer.Status() >= http.StatusInternalServerError:
			evt = logger.Error()
		default:
			evt = logger.Info()
		}

		logEvent := evt.
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", path).
			Str(domain.RequestIDKey, c.GetString(domain.RequestIDKey)).
			Str("ip", c.ClientIP())

		if currentUser, ok := domain.CurrentUserFromContext(c.Request.Context()); ok {
			logEvent = logEvent.Int64(domain.UserIDKey, currentUser.UserID).Str(domain.UserKey, currentUser.Account)
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
