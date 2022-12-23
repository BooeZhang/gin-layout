package config

// FeatureConfig 监控配置项。
type FeatureConfig struct {
	EnableProfiling bool `json:"profiling"      mapstructure:"profiling"`
	EnableMetrics   bool `json:"enable-metrics" mapstructure:"enable-metrics"`
}

// NewFeatureConfig 创建一个默认的监控配置项。
func NewFeatureConfig() *FeatureConfig {
	return &FeatureConfig{
		EnableMetrics:   true,
		EnableProfiling: true,
	}
}
