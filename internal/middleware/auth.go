package middleware

import (
	"github.com/gin-gonic/gin"

	"gin-layout/internal/reqctx"
	"gin-layout/internal/token"
)

func Auth(tokens token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawToken, err := token.ParseBearer(c.GetHeader("Authorization"))
		if err != nil {
			c.Error(err)
			return
		}

		claims, err := tokens.Parse(rawToken)
		if err != nil {
			c.Error(err)
			return
		}
		if claims.Type != token.TypeAccess {
			c.Error(token.ErrInvalidAccessToken)
			return
		}

		revoked, err := tokens.IsRevoked(c.Request.Context(), rawToken)
		if err != nil {
			c.Error(err)
			return
		}
		if revoked {
			c.Error(token.ErrTokenRevoked)
			return
		}

		contextUser := reqctx.CurrentUser{
			UserID:  claims.UserID,
			Account: claims.Subject,
		}

		c.Request = c.Request.WithContext(reqctx.WithCurrentUser(c.Request.Context(), contextUser))
		c.Request = c.Request.WithContext(reqctx.WithCurrentToken(c.Request.Context(), rawToken))

		c.Next()
	}
}
