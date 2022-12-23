package config

// ServerRunConfig 服务通用的选项
type ServerRunConfig struct {
	Mode        string   `json:"mode"        mapstructure:"mode"`
	Health      bool     `json:"health"     mapstructure:"health"`
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

// NewServerRunConfig 使用默认参数创建.
func NewServerRunConfig() *ServerRunConfig {
	return &ServerRunConfig{
		Mode:        "debug",
		Health:      true,
		Middlewares: []string{},
	}
}
