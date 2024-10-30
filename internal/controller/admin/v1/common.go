package v1

import (
	"github.com/gin-gonic/gin"

	"gin-layout/internal/model"
	srvv1 "gin-layout/internal/service/admin/v1"
	"gin-layout/pkg/response"
)

type CommonController struct {
	srv *srvv1.CommService
}

func NewCommonController() *CommonController {
	return &CommonController{srv: srvv1.NewCommService()}
}

// Login
// @Summary 登录
// @Schemes
// @Description 登录
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param   data body schema.LoginReq true "."
// @Success 200 {object} response.Response{data=model.LoginRes} "ok"
// @Router /login/ [post]
func (cc CommonController) Login(ctx *gin.Context) {
	var param model.LoginReq
	err := ctx.BindJSON(&param)
	if err != nil {
		response.Error(ctx, err, nil)
		return
	}
	data, err := cc.srv.Login(ctx, param.UserName, param.Password)
	if err != nil {
		response.Error(ctx, err, nil)
		return
	}
	response.Ok(ctx, nil, data)
}
