package config

import (
	"time"
)

// MySQL mysql配置项
type MySQL struct {
	Host                  string        `json:"host"                     mapstructure:"host"`
	Username              string        `json:"username"                 mapstructure:"username"`
	Password              string        `json:"password"                 mapstructure:"password"`
	Database              string        `json:"database"                 mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections"     mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `json:"max-open-connections"     mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time" mapstructure:"max-connection-life-time"`
	LogLevel              string        `json:"log-level"                mapstructure:"log-level"`

	SuperUser    string `json:"super-user" mapstructure:"super-user"`
	SuperUserPwd string `json:"super-user-pwd" mapstructure:"super-user-pwd"`
}
