package router

import (
	"encoding/base64"
	"github.com/BooeZhang/gin-layout/internal/apiserver/middleware"
	"strings"
	"time"

	"github.com/BooeZhang/gin-layout/internal/apiserver/datastore"
	"github.com/BooeZhang/gin-layout/internal/apiserver/model"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/response"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type loginInfo struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func newJWTAuth(store datastore.Factory) *jwt.GinJWTMiddleware {
	ginJwt, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            viper.GetString("jwt.Realm"),
		SigningAlgorithm: "HS256",
		Key:              []byte(viper.GetString("jwt.key")),
		Timeout:          viper.GetDuration("jwt.timeout") * time.Hour,
		MaxRefresh:       viper.GetDuration("jwt.max-refresh") * time.Hour,
		Authenticator:    authenticator(store),
		LoginResponse:    loginResponse(),
		LogoutResponse: func(c *gin.Context, code int) {
			response.Ok(c, nil, nil)
		},
		RefreshResponse: refreshResponse(),
		PayloadFunc:     payloadFunc(),
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims[jwt.IdentityKey]
		},
		IdentityKey:  middleware.UsernameKey,
		Authorizator: authorize(),
		Unauthorized: func(c *gin.Context, code int, message string) {
			response.Ok(c, nil, gin.H{
				"code": code,
				"msg":  message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		SendCookie:    true,
		TimeFunc:      time.Now,
		// TODO: HTTPStatusMessageFunc:
	})

	return ginJwt
}

// 校验登陆用户是否合法并生成jwt token
func authenticator(store datastore.Factory) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var login loginInfo
		var err error

		// support header and body both
		if c.Request.Header.Get("Authorization") != "" {
			login, err = parseWithHeader(c)
		} else {
			login, err = parseWithBody(c)
		}
		if err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		user, err := store.SysUser().GetSysUserByName(c.Request.Context(), login.Username)
		if err != nil {
			log.Errorf("get user information failed: %s", err.Error())

			return "", jwt.ErrFailedAuthentication
		}

		// Compare the login password with the user password.
		if err := user.Compare(login.Password); err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		return user, nil
	}
}

func loginResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		response.Ok(c, nil, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func refreshResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		response.Ok(c, nil, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		claims := jwt.MapClaims{}
		if u, ok := data.(model.SysUserModel); ok {
			claims[jwt.IdentityKey] = u.UserName
			claims["id"] = u.ID
			claims["username"] = u.UserName
			claims["is_super"] = u.IsSuper
		}

		return claims
	}
}

func parseWithHeader(c *gin.Context) (loginInfo, error) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		log.Errorf("get basic string from Authorization header failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	payload, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		log.Errorf("decode basic string: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		log.Errorf("parse payload failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}
	return loginInfo{
		Username: pair[0],
		Password: pair[1],
	}, nil
}

func parseWithBody(c *gin.Context) (loginInfo, error) {
	var login loginInfo
	if err := c.ShouldBindJSON(&login); err != nil {
		log.Errorf("parse login parameters: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	return login, nil
}

func authorize() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(string); ok {
			log.L(c).Infof("user `%s` is authenticated.", v)

			return true
		}

		return false
	}
}
