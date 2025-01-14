package v1

import (
	"github.com/gin-gonic/gin"

	"gin-layout/internal/model"
	"gin-layout/internal/repo"
	rp "gin-layout/internal/repo/mysql"
	"gin-layout/pkg/erroron"
	"gin-layout/pkg/jwtx"
	"gin-layout/store/mysqlx"
)

type CommService struct {
	userRepo repo.UserRepo
}

func NewCommService() *CommService {
	return &CommService{rp.NewUserRepo(mysqlx.GetDB())}
}

func (cs *CommService) Login(ctx *gin.Context, name, pwd string) (model.LoginRes, error) {
	var (
		res model.LoginRes
	)
	user, err := cs.userRepo.GetUserByName(ctx, name)
	if err != nil {
		return res, err
	}

	if user == nil || user.ID == 0 {
		return res, erroron.ErrNotFoundUser
	}

	if user.Compare(pwd) != nil {
		return res, erroron.ErrUserNameOrPwd
	}

	claims := jwtx.UserClaims{
		UserId:   user.ID,
		UserName: user.Account,
	}
	res.AccessToken, err = jwtx.GenAccessToken(claims)
	res.RefreshToken, err = jwtx.GenRefreshToken(claims)
	if err != nil {
		return res, err
	}
	return res, nil

}
