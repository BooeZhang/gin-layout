package jwtx

import (
	"github.com/BooeZhang/gin-layout/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type UserClaims struct {
	UserId   uint32
	UserName string
	Role     string
	RoleId   uint32
	Expire   time.Time
}

type jwtClaims struct {
	UserClaims
	jwt.RegisteredClaims
}

func GenToken(userClaims UserClaims) (string, error) {
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
		expire = cf.JwtConfig.Timeout * time.Hour
	}

	now := time.Now()
	expiresAt := time.Now().Add(expire)
	userClaims.Expire = expiresAt

	claims := &jwtClaims{
		UserClaims: userClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(now),       // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(expiresAt), // 过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func ParseToken(tokenString string) (*UserClaims, error) {
	var (
		key []byte
	)
	cf := config.GetConfig()
	if cf == nil {
		key = []byte("jwt")
	} else {
		key = []byte(cf.JwtConfig.Key)
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return key, nil
	})
	if token != nil {
		if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
			return &claims.UserClaims, nil
		}
	}
	return nil, err
}
