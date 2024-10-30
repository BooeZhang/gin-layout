package config

import (
	"time"
)

// Jwt JWT配置项.
type Jwt struct {
	Key            string        `json:"key"         mapstructure:"key"`
	AccessExpired  time.Duration `json:"access-expired"     mapstructure:"access-expired"`
	RefreshExpired time.Duration `json:"refresh-expired" mapstructure:"refresh-expired"`
}
