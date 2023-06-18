package config

// FeatureConfig 监控配置项。
type FeatureConfig struct {
	EnableProfiling bool `json:"profiling"      mapstructure:"profiling"`
	EnableMetrics   bool `json:"enable-metrics" mapstructure:"enable-metrics"`
}
