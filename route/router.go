package route

import (
	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/middleware"
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/response"
	"github.com/BooeZhang/gin-layout/store/mysql"
	"github.com/gin-gonic/gin"
)

func InitRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {
	jwtStrategy := middleware.NewJWTAuth(mysql.GetDB(), config.GetConfig().JwtConfig)
	errInit := jwtStrategy.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	g.POST("/login", jwtStrategy.LoginHandler)
	g.POST("/logout", jwtStrategy.LogoutHandler)
	g.POST("/refresh", jwtStrategy.RefreshHandler)

	g.NoRoute(jwtStrategy.MiddlewareFunc(), func(c *gin.Context) {
		response.Ok(c, erroron.ErrNotFound, nil)
	})

	return g
}
