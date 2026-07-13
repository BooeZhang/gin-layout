package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var cfg *Config

// Load 加载配置
func Load(configFile string) (*Config, error) {
	if configFile == "" {
		return nil, fmt.Errorf("please specify config file with -c")
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("toml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	dsn, err := BuildDatabaseDSN(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("build database DSN failed: %w", err)
	}
	cfg.Database.DSN = dsn

	return cfg, nil
}

// GetConfig 获取配置
func GetConfig() *Config {
	return cfg
}
