package config

// HttpServerConfig http服务配置项.
type HttpServerConfig struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port" mapstructure:"bind-port"`
	// ServerCert TLS 证书信息
	ServerCert CertKey `json:"tls"          mapstructure:"tls"`
}

// CertKey 证书相关配置
type CertKey struct {
	CertFile string `json:"cert-file"        mapstructure:"cert-file`
	KeyFile  string `json:"private-key-file" mapstructure:"private-key-file"`
}
