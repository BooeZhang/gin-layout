package v1

import (
	"github.com/BooeZhang/gin-layout/internal/apiserver/datastore"
	"github.com/BooeZhang/gin-layout/internal/apiserver/model"
	srvv1 "github.com/BooeZhang/gin-layout/internal/apiserver/service/v1"
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/response"
	"github.com/gin-gonic/gin"
)

type SysUserController struct {
	srv srvv1.Service
}

func NewSysUserController(store datastore.Factory) *SysUserController {
	return &SysUserController{srv: srvv1.NewService(store)}
}

// Create add new user to the storage.
func (su *SysUserController) Create(c *gin.Context) {
	log.L(c).Info("user create function called.")

	var r model.SysUserModel

	if err := c.ShouldBindJSON(&r); err != nil {
		response.Error(c, erroron.ErrParameter, nil)

		return
	}

	if err := su.srv.SysUser().Create(c, &r); err != nil {
		response.Error(c, err, nil)

		return
	}

	response.Ok(c, nil, "Ok")
}
