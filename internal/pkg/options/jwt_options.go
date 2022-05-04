package options

import (
	"time"
)

// JwtOptions JWT配置项.
type JwtOptions struct {
	Realm      string        `json:"realm"       mapstructure:"realm"`
	Key        string        `json:"key"         mapstructure:"key"`
	Timeout    time.Duration `json:"timeout"     mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"max-refresh" mapstructure:"max-refresh"`
}

// NewJwtOptions 创建一个默认的 JWT 配置项。
func NewJwtOptions() *JwtOptions {
	return &JwtOptions{
		Realm:      "jwt",
		Key:        "",
		Timeout:    1 * time.Hour,
		MaxRefresh: 1 * time.Hour,
	}
}
