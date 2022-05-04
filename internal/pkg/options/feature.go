package options

// FeatureOptions 监控配置项。
type FeatureOptions struct {
	EnableProfiling bool `json:"profiling"      mapstructure:"profiling"`
	EnableMetrics   bool `json:"enable-metrics" mapstructure:"enable-metrics"`
}

// NewFeatureOptions 创建一个默认的监控配置项。
func NewFeatureOptions() *FeatureOptions {
	return &FeatureOptions{
		EnableMetrics:   true,
		EnableProfiling: true,
	}
}
