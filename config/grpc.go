package config

import (
	"time"
)

// GRPC RPC配置项
type GRPC struct {
	Network           string        `json:"network" mapstructure:"network"`
	Addr              string        `json:"addr" mapstructure:"addr"`
	Timeout           time.Duration `json:"timeout" mapstructure:"timeout"`
	IdleTimeout       time.Duration `json:"idle-timeout" mapstructure:"idle-timeout"`
	MaxLifeTime       time.Duration `json:"max-life-time" mapstructure:"max-life-time"`
	ForceCloseWait    time.Duration `json:"force-close-wait" mapstructure:"force-close-wait"`
	KeepAliveInterval time.Duration `json:"keep-alive-interval" mapstructure:"keep-alive-interval"`
	KeepAliveTimeout  time.Duration `json:"keep-alive-timeout" mapstructure:"keep-alive-timeout"`
}
