package user

import (
	"github.com/BooeZhang/gin-layout/internal/service/v1/user"
	"github.com/BooeZhang/gin-layout/pkg/response"
	"github.com/gin-gonic/gin"
)

// Handler 用户业务handler
type Handler struct {
	svc user.Service
}

func NewUserHandler(_userSrv user.Service) *Handler {
	return &Handler{
		svc: _userSrv,
	}
}

// GetUserInfo
// @Summary 获取用户信息
// @Schemes
// @Description 获取用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /user/ [get]
func (uh *Handler) GetUserInfo(c *gin.Context) {
	uid := c.GetInt64("user_id")
	_user, err := uh.svc.GetById(c, uid)
	if err != nil {
		response.Error(c, err, nil)
	} else {
		response.Ok(c, nil, _user)
	}
}

func (uh *Handler) AddAuth(c *gin.Context) {

}
