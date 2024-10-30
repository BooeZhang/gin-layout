package config

type Log struct {
	Formatter string `json:"format" mapstructure:"format"`
	Level     string `json:"level" mapstructure:"level"`
}
