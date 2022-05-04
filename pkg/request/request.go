package request

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// GetUserID 获取用户id
func GetUserID(c *gin.Context) uint32 {
	claims := jwt.ExtractClaims(c)
	return claims["id"].(uint32)
}
