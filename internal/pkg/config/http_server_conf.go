package config

// HttpServerConfig http服务配置项.
type HttpServerConfig struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	// BindPort 设置 Listener 时被忽略，即使为 0 也会提供 HTTPS。
	BindPort int `json:"bind-port"    mapstructure:"bind-port"`
	// Required 设置为 true 意味着 BindPort 不能为零。
	Required bool
	// ServerCert TLS 证书信息
	ServerCert CertKey `json:"tls"          mapstructure:"tls"`
}

// CertKey 证书相关配置
type CertKey struct {
	CertFile string `json:"cert-file"        mapstructure:"cert-file"`
	KeyFile  string `json:"private-key-file" mapstructure:"private-key-file"`
}

// NewHttpServerConfig 创建一个带有默认参数的http服务.
func NewHttpServerConfig() *HttpServerConfig {
	return &HttpServerConfig{
		BindAddress: "0.0.0.0",
		BindPort:    8090,
		Required:    true,
		ServerCert: CertKey{
			CertFile: "",
			KeyFile:  "",
		},
	}
}
