package token

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"gin-layout/internal/reqctx"
)

const (
	TypeAccess  = "access"
	TypeRefresh = "refresh"
)

var (
	ErrInvalidAccessToken = errors.New("无效访问令牌")
	ErrUnauthenticated    = errors.New("未登录或非法访问")
	ErrTokenInvalid       = errors.New("token 无效")
	ErrTokenExpired       = errors.New("token 已过期")
	ErrTokenNotActive     = errors.New("token 不是活跃状态")
	ErrTokenRevoked       = errors.New("token 已失效")
)

type Pair struct {
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	UserID    int64
	Subject   string
	Type      string
	ExpiresAt time.Time
}

type Issuer interface {
	Issue(userID int64, subject string) (Pair, error)
	Parse(raw string) (*Claims, error)
}

type BlacklistRepository interface {
	Exists(ctx context.Context, tokenHash string) (bool, error)
	Add(ctx context.Context, tokenHash string, userID int64, expiresAt time.Time) error
}

type Manager interface {
	Issuer
	IsRevoked(ctx context.Context, raw string) (bool, error)
	Revoke(ctx context.Context, raw string, userID int64, expiresAt time.Time) error
	RevokeCurrent(ctx context.Context) (bool, error)
}

func ParseBearer(authHeader string) (string, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", ErrInvalidAccessToken
	}
	return strings.TrimPrefix(authHeader, prefix), nil
}

func tokenHash(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

type manager struct {
	issuer    Issuer
	blacklist BlacklistRepository
}

func NewManager(issuer Issuer, blacklist BlacklistRepository) Manager {
	return &manager{issuer: issuer, blacklist: blacklist}
}

func (m *manager) Issue(userID int64, subject string) (Pair, error) {
	return m.issuer.Issue(userID, subject)
}

func (m *manager) Parse(raw string) (*Claims, error) {
	return m.issuer.Parse(raw)
}

func (m *manager) IsRevoked(ctx context.Context, raw string) (bool, error) {
	return m.blacklist.Exists(ctx, tokenHash(raw))
}

func (m *manager) Revoke(ctx context.Context, raw string, userID int64, expiresAt time.Time) error {
	return m.blacklist.Add(ctx, tokenHash(raw), userID, expiresAt)
}

func (m *manager) RevokeCurrent(ctx context.Context) (bool, error) {
	user, ok := reqctx.CurrentUserFromContext(ctx)
	if !ok {
		return false, ErrUnauthenticated
	}

	raw, ok := reqctx.CurrentTokenFromContext(ctx)
	if !ok {
		return false, ErrUnauthenticated
	}

	claims, err := m.issuer.Parse(raw)
	if err != nil {
		return false, err
	}

	if err := m.Revoke(ctx, raw, user.UserID, claims.ExpiresAt); err != nil {
		return false, err
	}
	return true, nil
}
