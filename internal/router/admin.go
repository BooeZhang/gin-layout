package router

import (
	"github.com/gin-gonic/gin"

	"github.com/BooeZhang/gin-layout/core"
	admin2 "github.com/BooeZhang/gin-layout/internal/handler/v1/admin"
	"github.com/BooeZhang/gin-layout/pkg/auth"
	"github.com/BooeZhang/gin-layout/pkg/middleware"
)

// 管理后台路由
type admin struct {
}

var Admin = new(admin)

func (*admin) Load(g *gin.Engine) {
	st := core.GetStore()
	comm := admin2.NewCommHandler(st)
	// 登录
	g.POST("/v1/user/login", comm.Login)

	// user group
	ug := g.Group("/v1/user", middleware.JWTAuth(), middleware.NewAuthorizer(auth.GetEnforcer()))
	{
		user := admin2.NewUserHandler(st)

		ug.GET("", user.GetUserInfo)
		// ug.POST('/add_auth', ar.userHandler.)
	}
}
