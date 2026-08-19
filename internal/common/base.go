package common

import (
	"context"

	"github.com/rs/zerolog"

	"gin-layout/config"
	"gin-layout/internal/infra"
)

// BaseService 只收敛跨 service 的日志和管理员判定，避免每个 feature 重复实现。
type BaseService struct {
	logger *zerolog.Logger
	cfg    *config.Config
}

func NewBaseService(cfg *config.Config, logger *zerolog.Logger) BaseService {
	return BaseService{cfg: cfg, logger: logger}
}

func (s BaseService) Log(ctx context.Context) zerolog.Logger {
	if s.logger == nil {
		return zerolog.Nop()
	}
	return infra.LogFromContext(ctx, s.logger)
}

func (s BaseService) IsAdmin(account string) bool {
	return s.cfg != nil && account == s.cfg.Initializer.AdminAccount
}

func (s BaseService) IsAdminRole(roleCode string) bool {
	return s.cfg != nil && roleCode == s.cfg.Initializer.AdminRoleCode
}
