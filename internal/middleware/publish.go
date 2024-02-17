package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

// Publish publish a redisx event to specified redisx channel when some action occurred.
func Publish() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		var resource string

		pathSplit := strings.Split(c.Request.URL.Path, "/")
		if len(pathSplit) > 2 {
			resource = pathSplit[2]
		}

		method := c.Request.Method
		fmt.Println(method)

		switch resource {
		//case "policies":
		//	notify(method, load.NoticePolicyChanged)
		//case "secrets":
		//	notify(method, load.NoticeSecretChanged)
		default:
		}
	}
}
