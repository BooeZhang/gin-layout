package repo

import (
	"github.com/gin-gonic/gin"

	"gin-layout/internal/model"
)

// UserRepo 用户repo接口
type UserRepo interface {
	GetUserByName(ctx *gin.Context, account string) (*model.User, error)
	GetUserById(ctx *gin.Context, uid uint) (*model.User, error)
	GetUserByMobile(ctx *gin.Context, mobile string) (*model.User, error)
}
