package v1

import (
	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/model"
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/response"
	v1 "github.com/BooeZhang/gin-layout/service/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type SysUserController struct {
	srv v1.Service
}

func NewSysUserController(ms *gorm.DB, rs redis.UniversalClient, cf *config.Config) *SysUserController {
	return &SysUserController{srv: v1.NewService(ms, rs, cf)}
}

// Create add new user to the storage.
func (su *SysUserController) Create(c *gin.Context) {
	log.L(c).Info("user create function called.")

	var r model.SysUser

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
