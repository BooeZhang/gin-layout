package admin

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/BooeZhang/gin-layout/internal/repo"
	"github.com/BooeZhang/gin-layout/internal/repo/mysql"
	v1 "github.com/BooeZhang/gin-layout/internal/service/v1"
	"github.com/BooeZhang/gin-layout/pkg/constant"
	"github.com/BooeZhang/gin-layout/pkg/crypto/hash"
	"github.com/BooeZhang/gin-layout/pkg/erroron"
	"github.com/BooeZhang/gin-layout/pkg/jwtx"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/schema"
	"github.com/BooeZhang/gin-layout/store"
)

// CommService 用户业务逻辑
type CommService struct {
	ctx *v1.ServiceContext
	ur  repo.UserRepo
}

func NewCommService(s store.Storage) *CommService {
	return &CommService{
		ctx: v1.NewServiceContext(),
		ur:  mysql.NewUserRepo(s),
	}
}

// Login 登录
func (cs *CommService) Login(ctx context.Context, name, pwd string) (*schema.LoginRes, error) {
	var (
		res schema.LoginRes
	)
	user, err := cs.ur.GetUserByName(ctx, name)
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
		UserName: user.Account,
	}
	token, err := jwtx.GenToken(claims)
	if err != nil {
		log.L(ctx).Error("生成 token 失败", zap.Error(err))
		return nil, err

	}

	res.Token = hash.MD5String(token)
	key := constant.RedisKeyPrefixToken + res.Token
	cs.ctx.Rs.Set(ctx, key, token, time.Hour*cs.ctx.Cfg.JwtConfig.Timeout)
	return &res, nil
}
