package user

import (
	"context"
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/BooeZhang/gin-layout/pkg/jwtx"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/schema"
	"go.uber.org/zap"
)

// Login 登录
func (us *serviceImpl) Login(ctx context.Context, name, pwd string) (*schema.LoginRes, error) {
	var (
		res schema.LoginRes
	)
	user, err := us.ur.GetUserByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, erroron.ErrNotFoundUser
	}

	if user.Compare(pwd) != nil {
		return nil, erroron.ErrUserNameOrPwd
	}

	claims := jwtx.UserClaims{
		UserId:   user.ID,
		UserName: user.UserName,
	}
	res.Token, err = jwtx.GenToken(claims)
	if err != nil {
		log.L(ctx).Error("生成 token 失败", zap.Error(err))
		return nil, err

	}

	return &res, nil
}
