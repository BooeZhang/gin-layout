package config

import (
	"time"
)

// JwtConfig JWT配置项.
type JwtConfig struct {
	Realm      string        `json:"realm"       mapstructure:"realm"`
	Key        string        `json:"key"         mapstructure:"key"`
	Timeout    time.Duration `json:"timeout"     mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"max-refresh" mapstructure:"max-refresh"`
}

// NewJwtConfig 创建一个默认的 JWT 配置项。
func NewJwtConfig() *JwtConfig {
	return &JwtConfig{
		Realm:      "jwt",
		Key:        "",
		Timeout:    1 * time.Hour,
		MaxRefresh: 1 * time.Hour,
	}
}
