package router

import (
	"github.com/BooeZhang/gin-layout/internal/handler/v1/user"
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
	// login
	// user group
	ug := g.Group("/v1/user")
	{
		ug.GET("", ar.userHandler.GetUserInfo)
	}
}
