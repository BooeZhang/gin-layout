package user

import (
	"github.com/BooeZhang/gin-layout/internal/service/v1/user"
	"github.com/BooeZhang/gin-layout/pkg/response"
	"github.com/gin-gonic/gin"
)

// Handler 用户业务handler
type Handler struct {
	userSrv user.Service
}

func NewUserHandler(_userSrv user.Service) *Handler {
	return &Handler{
		userSrv: _userSrv,
	}
}

func (uh *Handler) GetUserInfo(c *gin.Context) {
	uid := c.GetInt64("user_id")
	user, err := uh.userSrv.GetById(c, uid)
	if err != nil {
		response.Error(c, err, nil)
	} else {
		response.Ok(c, nil, user)
	}
}
