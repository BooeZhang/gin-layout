package common

import (
	"gin-layout/config"
)

// BaseService 只收敛跨 service 的日志和管理员判定，避免每个 feature 重复实现。
type BaseService struct {
	cfg *config.Config
}

func NewBaseService(cfg *config.Config) BaseService {
	return BaseService{cfg: cfg}
}

func (s BaseService) IsAdmin(account string) bool {
	return s.cfg != nil && account == s.cfg.Initializer.AdminAccount
}

func (s BaseService) IsAdminRole(roleCode string) bool {
	return s.cfg != nil && roleCode == s.cfg.Initializer.AdminRoleCode
}
