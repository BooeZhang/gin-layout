package config

// ServerRunConfig 服务通用的选项
type ServerRunConfig struct {
	Debug       bool     `json:"debug"        mapstructure:"debug"`
	Health      bool     `json:"health"     mapstructure:"health"`
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

// NewServerRunConfig 使用默认参数创建.
func NewServerRunConfig() *ServerRunConfig {
	return &ServerRunConfig{
		Debug:       true,
		Health:      true,
		Middlewares: []string{},
	}
}
