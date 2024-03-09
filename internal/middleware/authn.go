package middleware

import (
	"github.com/BooeZhang/gin-layout/pkg/constant"
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/BooeZhang/gin-layout/pkg/jwtx"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/response"
	"github.com/BooeZhang/gin-layout/store/redisx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

// JWTAuth jwt 认证
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			response.Error(c, erroron.ErrNotLogin, nil)
			return
		}

		if !strings.HasPrefix(token, "Bearer ") && !strings.HasPrefix(token, "Basic") {
			response.Error(c, erroron.ErrTokenInvalid, nil)
			return
		}

		rs := redisx.GetRedis()
		key := constant.RedisKeyPrefixToken + token
		token, err := rs.Get(c, key).Result()
		if err != nil {
			log.L(c).Error("token 缓存获取失败", zap.Error(err))
			response.Error(c, erroron.ErrTokenInvalid, nil)
			return
		}

		claims, err := jwtx.ParseToken(token)
		if err != nil {
			log.L(c).Error("token 解析错误", zap.Error(err))
			response.Error(c, err, nil)
			return
		}

		c.Set("userClaims", claims)
		c.Set("Uid", claims.UserId)
		c.Next()
	}
}
