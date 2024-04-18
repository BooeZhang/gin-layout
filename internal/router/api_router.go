package router

import (
	"github.com/BooeZhang/gin-layout/internal/handler/v1/user"
	"github.com/BooeZhang/gin-layout/pkg/auth"
	middleware2 "github.com/BooeZhang/gin-layout/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// ApiRouter api 路由
type ApiRouter struct {
	userHandler *user.Handler
}

func NewApiRouter(userHandler *user.Handler) *ApiRouter {
	return &ApiRouter{userHandler: userHandler}
}

func (ar *ApiRouter) Load(g *gin.Engine) {
	// 登录
	g.POST("/v1/user/login", ar.userHandler.Login)

	// user group
	ug := g.Group("/v1/user", middleware2.JWTAuth(), middleware2.NewAuthorizer(auth.GetEnforcer()))
	{
		ug.GET("", ar.userHandler.GetUserInfo)
		// ug.POST('/add_auth', ar.userHandler.)
	}
}
