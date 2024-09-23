package config

import (
	"time"
)

// Jwt JWT配置项.
type Jwt struct {
	Realm      string        `json:"realm"       mapstructure:"realm"`
	Key        string        `json:"key"         mapstructure:"key"`
	Timeout    time.Duration `json:"timeout"     mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"max-refresh" mapstructure:"max-refresh"`
}
