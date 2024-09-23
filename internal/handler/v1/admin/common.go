package admin

import (
	"github.com/gin-gonic/gin"

	srvv1 "github.com/BooeZhang/gin-layout/internal/service/v1/admin"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/response"
	"github.com/BooeZhang/gin-layout/pkg/schema"
	"github.com/BooeZhang/gin-layout/store"
)

// CommHandler 公共业务 handler
type CommHandler struct {
	srv *srvv1.CommService
}

func NewCommHandler(s store.Storage) *CommHandler {
	return &CommHandler{srv: srvv1.NewCommService(s)}
}

// Login
// @Summary 登录
// @Schemes
// @Description 登录
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param   data body schema.LoginReq true "."
// @Success 200 {object} response.Response{data=schema.LoginRes} "ok"
// @Router /user/login/ [post]
func (ch CommHandler) Login(c *gin.Context) {
	var param schema.LoginReq

	err := c.ShouldBindJSON(&param)
	if err != nil {
		log.L(c).Error(err.Error())
		response.Error(c, err, nil)
		return
	}
	data, err := ch.srv.Login(c, param.Username, param.Password)
	if err != nil {
		response.Error(c, err, nil)
		return
	}
	response.Ok(c, nil, data)
}
