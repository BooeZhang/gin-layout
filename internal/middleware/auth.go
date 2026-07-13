package middleware

import (
	"github.com/gin-gonic/gin"

	"gin-layout/internal/common"
	"gin-layout/internal/domain"
)

func Auth(tokens common.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawToken, err := common.ParseBearer(c.GetHeader("Authorization"))
		if err != nil {
			c.Error(err)
			return
		}

		claims, err := tokens.Parse(rawToken)
		if err != nil {
			c.Error(err)
			return
		}

		revoked, err := tokens.IsRevoked(c.Request.Context(), rawToken)
		if err != nil {
			c.Error(err)
			return
		}
		if revoked {
			c.Error(domain.ErrTokenRevoked)
			return
		}

		contextUser := domain.CurrentUser{
			UserID:  claims.UserID,
			Account: claims.Subject,
		}

		c.Request = c.Request.WithContext(domain.WithCurrentUser(c.Request.Context(), contextUser))
		c.Request = c.Request.WithContext(domain.WithCurrentToken(c.Request.Context(), rawToken))

		c.Next()
	}
}
