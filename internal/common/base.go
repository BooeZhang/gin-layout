package common

import (
	"context"
	"os"

	"github.com/rs/zerolog"

	"gin-layout/config"
	"gin-layout/internal/domain"
)

type BaseService struct {
	Logger *zerolog.Logger
	Cfg    *config.Config
}

func NewBaseService(cfg *config.Config, logger *zerolog.Logger) BaseService {
	if logger == nil {
		l := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger = &l
	}
	return BaseService{Logger: logger, Cfg: cfg}
}

func (s BaseService) Log(ctx context.Context) zerolog.Logger {
	logCtx := s.Logger.With()
	if requestID, ok := domain.RequestIDFromContext(ctx); ok {
		logCtx = logCtx.Str(domain.RequestIDKey, requestID)
	}
	if user, ok := domain.CurrentUserFromContext(ctx); ok {
		logCtx = logCtx.Int64(domain.UserIDKey, user.UserID).Str(domain.UserKey, user.Account)
	}
	return logCtx.Logger()
}

func (s BaseService) IsAdmin(account string) bool {
	return account == s.Cfg.Initializer.AdminAccount
}

func (s BaseService) IsAdminRole(roleCode string) bool {
	return roleCode == s.Cfg.Initializer.AdminRoleCode
}
