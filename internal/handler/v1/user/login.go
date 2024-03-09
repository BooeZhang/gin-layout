package user

import (
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
// @Success 200 body schema.LoginRes()
// @Success 200 {object} response.Response{data=schema.LoginRes} "ok"
// @Router /user/login/ [post]
func (uh *Handler) Login(c *gin.Context) {
	var param schema.LoginReq

	err := c.ShouldBindJSON(&param)
	if err != nil {
		log.L(c).Error(err.Error())
		response.Error(c, err, nil)
		return
	}
	data, err := uh.svc.Login(c, param.Username, param.Password)
	if err != nil {
		response.Error(c, err, nil)
		return
	}
	response.Ok(c, nil, data)
}
