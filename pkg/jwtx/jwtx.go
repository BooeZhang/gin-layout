package jwtx

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"gin-layout/config"
	"gin-layout/pkg/crypto/aes"
	"gin-layout/pkg/erroron"
)

var (
	aesKey = []byte("kFmjRMHBIlLH8I91s523MM6LGh3NL3pp")
)

type UserClaims struct {
	jwt.Claims
	UserId   uint   `json:"user_id"`
	UserName string `json:"user_name"`
}

type jwtClaims struct {
	UserClaims
	jwt.RegisteredClaims
}

// GenAccessToken 生成访问 token
func GenAccessToken(claims UserClaims) (string, error) {
	var (
		key    []byte
		expire time.Duration
	)
	cf := config.GetConfig()
	if cf == nil {
		key = []byte("jwt")
		expire = 1 * time.Hour
	} else {
		key = []byte(cf.JwtConfig.Key)
		expire = cf.JwtConfig.AccessExpired * time.Minute
	}

	now := time.Now()
	expiresAt := time.Now().Add(expire)
	_jwtClaims := &jwtClaims{
		UserClaims: claims,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(now),       // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(expiresAt), // 过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, _jwtClaims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	tokenString, err = aes.EncryptToBase64([]byte(tokenString), aesKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GenRefreshToken 生成刷新 token
func GenRefreshToken(claims UserClaims) (string, error) {
	var (
		key    []byte
		expire time.Duration
	)
	cf := config.GetConfig()
	if cf == nil {
		key = []byte("jwt")
		expire = 3 * time.Hour
	} else {
		key = []byte(cf.JwtConfig.Key)
		expire = cf.JwtConfig.RefreshExpired * time.Hour
	}

	now := time.Now()
	expiresAt := time.Now().Add(expire)
	_jwtClaims := &jwtClaims{
		UserClaims: claims,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(now),       // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(expiresAt), // 过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, _jwtClaims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	tokenString, err = aes.EncryptToBase64([]byte(tokenString), aesKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken token 校验
func ParseToken(token string) (*UserClaims, error) {
	var (
		key []byte
	)
	cf := config.GetConfig()
	if cf == nil {
		key = []byte("jwt")
	} else {
		key = []byte(cf.JwtConfig.Key)
	}

	tokenByte, err := aes.DecryptFromBase64(token, aesKey)
	if err != nil {
		log.Error().Err(err).Str("token", token).Msg("token error")
		return nil, err
	}

	_token, err := jwt.ParseWithClaims(string(tokenByte), &jwtClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return key, nil
	})

	if err != nil {
		log.Error().Err(err).Str("token", token).Msg("token error")
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, erroron.ErrTokenExpired
		default:
			return nil, erroron.ErrTokenInvalid
		}
	}

	if _token != nil && _token.Valid {
		if claims, ok := _token.Claims.(*jwtClaims); ok {
			return &claims.UserClaims, nil
		}
	}
	return nil, err
}

// GetToken 获取 token
func GetToken(ctx *gin.Context) (string, error) {
	tokenStr := ctx.Request.Header.Get("Authorization")
	if tokenStr == "" {
		return "", errors.New("token not obtained")
	}

	if !strings.HasPrefix(tokenStr, "Bearer ") {
		return "", errors.New("token format error")
	}

	splitToken := strings.Split(tokenStr, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("token format error")
	}

	return strings.TrimSpace(splitToken[1]), nil
}

// GetClaims 获取 token Claims
func GetClaims(ctx *gin.Context) *UserClaims {
	token, err := GetToken(ctx)
	if err != nil {
		return &UserClaims{}
	}
	claims, err := ParseToken(token)
	if err != nil {
		return &UserClaims{}
	}

	return claims
}

// GetUserID Get user id
func GetUserID(ctx *gin.Context) uint {
	return GetClaims(ctx).UserId
}

// GetUserName 获取用户名
func GetUserName(ctx *gin.Context) string {
	return GetClaims(ctx).UserName
}
