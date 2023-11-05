package config

// CasbinConf casbin 配置
type CasbinConf struct {
	ModelPath string `json:"model_path" mapstructure:"model_path"`
}
