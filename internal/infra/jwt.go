package infra

import (
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"

	"gin-layout/config"
	"gin-layout/internal/token"
)

type jwtClaims struct {
	UserID int64
	Type   string
	jwtlib.RegisteredClaims
}

type JWTIssuer struct {
	secret         []byte
	accessExpired  time.Duration
	refreshExpired time.Duration
}

func NewJWTIssuer(cfg *config.JWTConfig) *JWTIssuer {
	return &JWTIssuer{
		secret:         []byte(cfg.Secret),
		accessExpired:  time.Duration(cfg.AccessExpired) * time.Minute,
		refreshExpired: time.Duration(cfg.RefreshExpired) * time.Minute,
	}
}

func (i *JWTIssuer) Issue(userID int64, subject string) (token.Pair, error) {
	accessToken, err := i.sign(userID, subject, token.TypeAccess, i.accessExpired)
	if err != nil {
		return token.Pair{}, err
	}
	refreshToken, err := i.sign(userID, subject, token.TypeRefresh, i.refreshExpired)
	if err != nil {
		return token.Pair{}, err
	}
	return token.Pair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (i *JWTIssuer) Parse(raw string) (*token.Claims, error) {
	parsed := &jwtClaims{}
	jwtToken, err := jwtlib.ParseWithClaims(raw, parsed, func(jwtToken *jwtlib.Token) (any, error) {
		if jwtToken.Method != jwtlib.SigningMethodHS256 {
			return nil, fmt.Errorf("%w: %s", token.ErrTokenInvalid, jwtToken.Header["alg"])
		}
		return i.secret, nil
	}, jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, normalizeJWTError(err)
	}
	if !jwtToken.Valid {
		return nil, token.ErrTokenInvalid
	}

	var expiresAt time.Time
	if parsed.ExpiresAt != nil {
		expiresAt = parsed.ExpiresAt.Time
	}
	return &token.Claims{
		UserID:    parsed.UserID,
		Subject:   parsed.Subject,
		Type:      parsed.Type,
		ExpiresAt: expiresAt,
	}, nil
}

func (i *JWTIssuer) sign(userID int64, subject, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	jwtToken := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, jwtClaims{
		UserID: userID,
		Type:   tokenType,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwtlib.NewNumericDate(now),
			Subject:   subject,
		},
	})
	return jwtToken.SignedString(i.secret)
}

func normalizeJWTError(err error) error {
	switch {
	case errors.Is(err, jwtlib.ErrSignatureInvalid):
		return token.ErrTokenInvalid
	case errors.Is(err, jwtlib.ErrTokenExpired):
		return token.ErrTokenExpired
	case errors.Is(err, jwtlib.ErrTokenNotValidYet):
		return token.ErrTokenNotActive
	default:
		return token.ErrTokenInvalid
	}
}
