package router

import (
	v1 "github.com/BooeZhang/gin-layout/internal/apiserver/controller/v1"
	"github.com/BooeZhang/gin-layout/internal/apiserver/datastore/mysql"
	"github.com/BooeZhang/gin-layout/internal/pkg/options"
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/response"
	"github.com/gin-gonic/gin"
)

func InitRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {
	opts := options.GetOptions()
	storeIns, _ := mysql.GetMysqlFactoryOr(opts.MySQLOptions)
	//storeCache, _ := redis.GetRedisFactoryOr(opts.RedisOptions)
	jwtStrategy := newJWTAuth(storeIns)
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
	{
		sysuserController := v1.NewSysUserContorller(storeIns)
		g.POST("add", sysuserController.Create)
	}

	return g
}
