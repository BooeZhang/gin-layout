package router

import (
	"github.com/gin-gonic/gin"

	v1 "gin-layout/internal/controller/admin/v1"
	"gin-layout/middleware"
)

type _admin struct{}

var Admin = _admin{}

func (_admin) Load(r *gin.Engine) {
	comm := v1.NewCommonController()
	g := r.Group("admin/v1/")
	g.POST("login", comm.Login)
	g.Use(middleware.Authn())

	user := g.Group("/user")
	{
		u := v1.NewUserController()
		user.GET("user_info", u.UserInfo)
	}
}
