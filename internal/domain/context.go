package domain

import (
	"context"
	"time"
)

type CurrentUser struct {
	UserID  int64
	Account string
}

type contextKey string

const (
	CtxUserKey    contextKey = "currentUser"
	CtxTokenKey   contextKey = "currentToken"
	CtxRequestID  contextKey = "requestID"
	RequestIDKey             = "requestID"
	UserIDKey                = "userID"
	UserKey                  = "user"
	RequestIDHeader          = "X-Request-ID"
)

func WithCurrentUser(ctx context.Context, user CurrentUser) context.Context {
	return context.WithValue(ctx, CtxUserKey, user)
}

func CurrentUserFromContext(ctx context.Context) (CurrentUser, bool) {
	u, ok := ctx.Value(CtxUserKey).(CurrentUser)
	return u, ok
}

func WithCurrentToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, CtxTokenKey, token)
}

func CurrentTokenFromContext(ctx context.Context) (string, bool) {
	t, ok := ctx.Value(CtxTokenKey).(string)
	return t, ok
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, CtxRequestID, requestID)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	r, ok := ctx.Value(CtxRequestID).(string)
	return r, ok
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenClaims struct {
	UserID  int64
	Subject string
	ExpiresAt time.Time
}

