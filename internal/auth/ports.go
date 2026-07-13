package auth

import (
	"context"
	"time"

	"gin-layout/internal/domain"
)

type UserProvider interface {
	FindByAccount(ctx context.Context, account string) (*domain.User, error)
	UpdateLastLogin(ctx context.Context, userID int64, at time.Time) error
}
