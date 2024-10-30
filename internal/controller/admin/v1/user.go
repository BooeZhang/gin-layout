package v1

import (
	"github.com/gin-gonic/gin"

	srvv1 "gin-layout/internal/service/admin/v1"
	"gin-layout/pkg/request"
	"gin-layout/pkg/response"
)

type UserController struct {
	srv *srvv1.UserService
}

func NewUserController() *UserController {
	return &UserController{srv: srvv1.NewUserService()}
}

// UserInfo
// @Summary 用户信息
// @Schemes
// @Description 用户信息
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param   data query request.ByID true "."
// @Success 200 {object} response.Response{data=model.User} "ok"
// @Router /login/ [post]
func (uc UserController) UserInfo(ctx *gin.Context) {
	var param request.ByID
	err := ctx.BindQuery(&param)
	if err != nil {
		response.Error(ctx, err, nil)
		return
	}
	data, err := uc.srv.UserInfo(ctx, param.ID)
	if err != nil {
		response.Error(ctx, err, nil)
		return
	}
	response.Ok(ctx, nil, data)
}
