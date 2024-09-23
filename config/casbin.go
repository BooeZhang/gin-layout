package config

// Casbin casbin 配置
type Casbin struct {
	ModelPath string `json:"model_path" mapstructure:"model_path"`
}
