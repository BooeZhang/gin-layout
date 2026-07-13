package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"gin-layout/internal/domain"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(domain.RequestIDHeader)
		if requestID == "" {
			requestID = newRequestID()
			c.Request.Header.Set(domain.RequestIDHeader, requestID)
		}

		c.Set(domain.RequestIDKey, requestID)
		c.Header(domain.RequestIDHeader, requestID)
		c.Request = c.Request.WithContext(domain.WithRequestID(c.Request.Context(), requestID))

		c.Next()
	}
}

func newRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return cast.ToString(time.Now().UnixNano())
	}
	return hex.EncodeToString(b[:])
}
