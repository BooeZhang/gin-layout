package user

import (
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/response"
	"github.com/BooeZhang/gin-layout/pkg/schema"
	"github.com/gin-gonic/gin"
)

// Login
// @Summary 登录
// @Schemes
// @Description 登录
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param   data body schema.LoginReq true "."
// @Success 200 body schema.LoginRes
// @Router /user/login/ [post]
func (uh *Handler) Login(c *gin.Context) {
	var param schema.LoginReq

	err := c.ShouldBindJSON(&param)
	if err != nil {
		log.L(c).Error(err.Error())
		response.Error(c, erroron.ErrParameter, nil)
		return
	}
	data, err := uh.userSrv.Login(c, param.Username, param.Password)
	if err != nil {
		response.Error(c, err, nil)
		return
	}
	response.Ok(c, nil, data)
}
