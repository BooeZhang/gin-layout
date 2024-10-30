package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"gin-layout/pkg/erroron"
	"gin-layout/pkg/jwtx"
	"gin-layout/pkg/response"
)

func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwtx.GetToken(c)
		if err != nil {
			log.Error().Err(err).Str("token", token).Msg("token错误")
			response.Error(c, erroron.ErrTokenInvalid, nil)
			return
		}

		claims, err := jwtx.ParseToken(token)
		if err != nil {
			log.Error().Err(err).Str("token", token).Msg("token 解析错误")
			response.Error(c, err, nil)
			return
		}

		c.Set("claims", claims)
		c.Set("user", claims.UserName)
		c.Next()
	}
}
