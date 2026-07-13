package common

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"gin-layout/internal/domain"
)

// ParseBearer extracts the Bearer token from an Authorization header.
func ParseBearer(authHeader string) (string, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", domain.ErrInvalidAccessToken
	}
	return strings.TrimPrefix(authHeader, prefix), nil
}

// TokenHash returns the SHA-256 hex digest of the raw token.
func TokenHash(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

type tokenManager struct {
	issuer    TokenIssuer
	blacklist TokenBlacklistRepository
}

func NewTokenManager(issuer TokenIssuer, blacklist TokenBlacklistRepository) TokenManager {
	return &tokenManager{issuer: issuer, blacklist: blacklist}
}

func (m *tokenManager) Issue(userID int64, subject string) (domain.TokenPair, error) {
	return m.issuer.Issue(userID, subject)
}

func (m *tokenManager) Parse(raw string) (*domain.TokenClaims, error) {
	return m.issuer.Parse(raw)
}

func (m *tokenManager) IsRevoked(ctx context.Context, raw string) (bool, error) {
	return m.blacklist.Exists(ctx, TokenHash(raw))
}

func (m *tokenManager) Revoke(ctx context.Context, raw string, userID int64, expiresAt time.Time) error {
	return m.blacklist.Add(ctx, TokenHash(raw), userID, expiresAt)
}

func (m *tokenManager) RevokeCurrent(ctx context.Context) (bool, error) {
	user, ok := domain.CurrentUserFromContext(ctx)
	if !ok {
		return false, domain.ErrNotLogin
	}

	raw, ok := domain.CurrentTokenFromContext(ctx)
	if !ok {
		return false, domain.ErrNotLogin
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
