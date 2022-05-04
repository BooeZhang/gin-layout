package options

// ServerRunOptions 服务通用的选项
type ServerRunOptions struct {
	Mode        string   `json:"mode"        mapstructure:"mode"`
	Healthz     bool     `json:"healthz"     mapstructure:"healthz"`
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

// NewServerRunOptions 使用默认参数创建.
func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		Mode:        "debug",
		Healthz:     true,
		Middlewares: []string{},
	}
}
