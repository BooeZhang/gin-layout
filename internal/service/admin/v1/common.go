package v1

import (
	"github.com/gin-gonic/gin"

	"gin-layout/internal/model"
)

type CommService struct {
	us *UserService
}

func NewCommService() *CommService {
	return &CommService{us: NewUserService()}
}

func (cs *CommService) Login(ctx *gin.Context, name, pwd string) (model.LoginRes, error) {
	return cs.us.Login(ctx, name, pwd)
}
