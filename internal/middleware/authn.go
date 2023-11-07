package middleware

import (
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/BooeZhang/gin-layout/pkg/jwtx"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWTAuth jwt 认证
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			response.Error(c, erroron.ErrNotLogin, nil)
			return
		}
		claims, err := jwtx.ParseToken(token)
		if err != nil {
			log.L(c).Error("token 解析错误", zap.Error(err))
			response.Error(c, erroron.ErrTokenInvalid, nil)
			return
		}

		c.Set("userClaims", claims)
		c.Next()
	}
}
