package log

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type Format string

const (
	FormatJSON    Format = "json"
	FormatConsole Format = "console"
)

type Config struct {
	Level       string       `mapstructure:"level"`
	Format      Format       `mapstructure:"format"`
	Service     string       `mapstructure:"service"`
	Environment string       `mapstructure:"environment"`
	Version     string       `mapstructure:"version"`
	Caller      bool         `mapstructure:"caller"`
	Color       bool         `mapstructure:"color"`
	TimeFormat  string       `mapstructure:"time_format"`
	Output      OutputConfig `mapstructure:"output"`
}

type OutputConfig struct {
	Stdout bool        `mapstructure:"stdout"`
	Stderr bool        `mapstructure:"stderr"`
	File   *FileConfig `mapstructure:"file"`
}

type FileConfig struct {
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
	LocalTime  bool   `mapstructure:"local_time"`
}

func Development(service string) Config {
	return Config{
		Level: "debug", Format: FormatConsole, Service: service,
		Environment: "development", Color: true, Caller: true,
		Output: OutputConfig{Stdout: true},
	}
}

func Production(service string) Config {
	return Config{
		Level: "info", Format: FormatJSON, Service: service,
		Environment: "production", Caller: true,
		Output: OutputConfig{Stdout: true},
	}
}

func (c Config) validate() (Config, error) {
	if c.TimeFormat == "" {
		c.TimeFormat = time.RFC3339Nano
	}

	switch c.Format {
	case FormatJSON, FormatConsole:
	default:
		return Config{}, fmt.Errorf("unsupported log format %q", c.Format)
	}

	if c.Level == "" {
		return Config{}, fmt.Errorf("log level must not be empty")
	}
	if _, err := zerolog.ParseLevel(c.Level); err != nil {
		return Config{}, fmt.Errorf("invalid log level %q: %w", c.Level, err)
	}

	if !c.Output.Stdout && !c.Output.Stderr && c.Output.File == nil {
		return Config{}, fmt.Errorf("at least one log output must be configured")
	}

	if c.Output.File != nil {
		f := c.Output.File
		if f.Filename == "" {
			return Config{}, fmt.Errorf("log file filename must not be empty")
		}
		if f.MaxSize < 0 {
			return Config{}, fmt.Errorf("log file max_size must be >= 0")
		}
		if f.MaxBackups < 0 {
			return Config{}, fmt.Errorf("log file max_backups must be >= 0")
		}
		if f.MaxAge < 0 {
			return Config{}, fmt.Errorf("log file max_age must be >= 0")
		}
	}

	return c, nil
}
