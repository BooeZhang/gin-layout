package config

// HttpServer http服务配置项.
type HttpServer struct {
	BindAddress     string   `json:"bind-address" mapstructure:"bind-address"`
	BindPort        int      `json:"bind-port" mapstructure:"bind-port"`
	Debug           bool     `json:"debug"       mapstructure:"debug"`
	Health          bool     `json:"health"      mapstructure:"health"`
	Middlewares     []string `json:"middlewares" mapstructure:"middlewares"`
	EnableMetrics   bool     `json:"enable-metrics" mapstructure:"enable-metrics"`
	EnableProfiling bool     `json:"enable-profiling" mapstructure:"enable-profiling"`
	ServerCert      CertKey  `json:"tls"          mapstructure:"tls"` // ServerCert TLS 证书信息
	CasbinModelPath string   `json:"casbin-model-path" mapstructure:"casbin-model-path"`
}

// CertKey 证书相关配置
type CertKey struct {
	CertFile string `json:"cert-file"        mapstructure:"cert-file"`
	KeyFile  string `json:"private-key-file" mapstructure:"private-key-file"`
}
