package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gin-layout/internal/model"
	"gin-layout/internal/repo"
	rp "gin-layout/internal/repo/mysql"
	"gin-layout/pkg/erroron"
	"gin-layout/pkg/jwtx"
	"gin-layout/store/mysqlx"
)

type UserService struct {
	userRepo repo.UserRepo
}

func NewUserService() *UserService {
	return &UserService{rp.NewUserRepo(mysqlx.GetDB())}
}

func (us UserService) Login(ctx *gin.Context, name, pwd string) (model.LoginRes, error) {
	var (
		res model.LoginRes
	)
	user, err := us.userRepo.GetUserByName(ctx, name)
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
	if err != nil {
		return res, err
	}
	res.RefreshToken, err = jwtx.GenRefreshToken(claims)
	if err != nil {
		return res, err
	}
	return res, nil

}

func (us UserService) UserInfo(ctx *gin.Context, ID uint) (*model.User, error) {
	if ID == 0 {
		ID = jwtx.GetUserID(ctx)
	}
	user, err := us.userRepo.GetUserById(ctx, ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, erroron.ErrNotFoundUser
	}

	return user, nil

}
