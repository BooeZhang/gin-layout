package config

import (
	"encoding/json"
	"fmt"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const configFlagName = "config"

var (
	c       *Config
	cfgFile string
)

func init() {
	pflag.StringVarP(&cfgFile, "config", "c", cfgFile, "Read configuration from specified `FILE`, "+
		"support JSON, TOML, YAML, HCL, or Java properties formats.")
}

// Config 配置选项
type Config struct {
	ServerRunConfig  *ServerRunConfig  `json:"server"   mapstructure:"server"`
	GRPCConfig       *GRPCConfig       `json:"grpc"     mapstructure:"grpc"`
	HttpServerConfig *HttpServerConfig `json:"http"     mapstructure:"http"`
	MySQLConfig      *MySQLConfig      `json:"mysql"    mapstructure:"mysql"`
	RedisConfig      *RedisConfig      `json:"redis"    mapstructure:"redis"`
	JwtConfig        *JwtConfig        `json:"jwt"      mapstructure:"jwt"`
	LogConfig        *log.Options      `json:"log"      mapstructure:"log"`
	FeatureConfig    *FeatureConfig    `json:"feature"  mapstructure:"feature"`
}

// DefaultConfig 默认配置选项
func DefaultConfig() *Config {
	c = &Config{
		ServerRunConfig:  NewServerRunConfig(),
		GRPCConfig:       NewGRPCConfig(),
		HttpServerConfig: NewHttpServerConfig(),
		MySQLConfig:      NewMySQLConfig(),
		RedisConfig:      NewRedisConfig(),
		JwtConfig:        NewJwtConfig(),
		LogConfig:        log.NewOptions(),
		FeatureConfig:    NewFeatureConfig(),
	}

	return c
}

// AddConfigFlag 读取配置
func AddConfigFlag(basename string, fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(configFlagName))

	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(basename), "-", "_", -1))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			viper.AddConfigPath(".")
			viper.AddConfigPath("etc/")
			viper.SetConfigName(basename)
		}

		if err := viper.ReadInConfig(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration file(%s): %v\n", cfgFile, err)
			os.Exit(1)
		}

		fmt.Println(viper.AllSettings())
		if err := viper.Unmarshal(&c); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: Unable to decode into struct file(%s): %v\n", cfgFile, err)
			os.Exit(1)
		}
	})
}

// String 配置字符输出
func (o *Config) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

func GetConfig() *Config {
	return c
}
