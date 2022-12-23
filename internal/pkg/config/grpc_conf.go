package config

// GRPCConfig RPC配置项
type GRPCConfig struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
	MaxMsgSize  int    `json:"max-msg-size" mapstructure:"max-msg-size"`
}

// NewGRPCConfig  创建默认的 RPC配置项
func NewGRPCConfig() *GRPCConfig {
	return &GRPCConfig{
		BindAddress: "0.0.0.0",
		BindPort:    8081,
		MaxMsgSize:  4 * 1024 * 1024,
	}
}
