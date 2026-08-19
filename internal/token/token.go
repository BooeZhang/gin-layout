package token

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"gin-layout/internal/domain"
	"gin-layout/internal/reqctx"
)

const (
	TypeAccess  = "access"
	TypeRefresh = "refresh"
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

var (
	ErrInvalidAccessToken = domain.NewDomainError(50010, http.StatusUnauthorized, "无效访问令牌")
	ErrNotLogin           = domain.NewDomainError(50011, http.StatusUnauthorized, "未登录或非法访问")
	ErrTokenInvalid       = domain.NewDomainError(50012, http.StatusUnauthorized, "token 无效")
	ErrTokenExpired       = domain.NewDomainError(50060, http.StatusUnauthorized, "token 已过期")
	ErrTokenNotActive     = domain.NewDomainError(50070, http.StatusUnauthorized, "token 不是活跃状态")
	ErrTokenRevoked       = domain.NewDomainError(50061, http.StatusUnauthorized, "token 已失效")
)

func ParseBearer(authHeader string) (string, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", ErrInvalidAccessToken
	}
	return strings.TrimPrefix(authHeader, prefix), nil
}

func TokenHash(raw string) string {
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
	return m.blacklist.Exists(ctx, TokenHash(raw))
}

func (m *manager) Revoke(ctx context.Context, raw string, userID int64, expiresAt time.Time) error {
	return m.blacklist.Add(ctx, TokenHash(raw), userID, expiresAt)
}

func (m *manager) RevokeCurrent(ctx context.Context) (bool, error) {
	user, ok := reqctx.CurrentUserFromContext(ctx)
	if !ok {
		return false, ErrNotLogin
	}

	raw, ok := reqctx.CurrentTokenFromContext(ctx)
	if !ok {
		return false, ErrNotLogin
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
