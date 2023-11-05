package router

import (
	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/internal/handler/v1/user"
	"github.com/BooeZhang/gin-layout/internal/middleware"
	"github.com/BooeZhang/gin-layout/pkg/auth"
	"github.com/BooeZhang/gin-layout/store/mysql"
	"github.com/gin-contrib/authz"
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
	db := mysql.GetDB()
	cf := config.GetConfig()
	jwtMiddleware := middleware.NewJWT(db, cf.JwtConfig)
	g.POST("/login", jwtMiddleware.LoginHandler)

	// user group
	ug := g.Group("/v1/user", authz.NewAuthorizer(auth.GetEnforcer()))
	{
		ug.GET("", ar.userHandler.GetUserInfo)
		//ug.POST('/add_auth', ar.userHandler.)
	}
}
