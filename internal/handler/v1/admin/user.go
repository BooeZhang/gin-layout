package admin

import (
	"github.com/gin-gonic/gin"

	srvv1 "github.com/BooeZhang/gin-layout/internal/service/v1/admin"
	"github.com/BooeZhang/gin-layout/pkg/response"
	"github.com/BooeZhang/gin-layout/store"
)

// UserHandler 用户业务 handler
type UserHandler struct {
	srv *srvv1.UserService
}

// NewUserHandler 初始化用户业务 handler 对象
func NewUserHandler(s store.Storage) *UserHandler {
	return &UserHandler{
		srv: srvv1.NewUserService(s),
	}
}

// GetUserInfo
// @Summary 获取用户信息
// @Schemes
// @Description 获取用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {string}
// @Router /user/ [get]
func (uh *UserHandler) GetUserInfo(c *gin.Context) {
	uid := c.GetInt64("user_id")
	_user, err := uh.srv.GetById(c, uid)
	if err != nil {
		response.Error(c, err, nil)
	} else {
		response.Ok(c, nil, _user)
	}
}

func (uh *UserHandler) AddAuth(c *gin.Context) {

}
