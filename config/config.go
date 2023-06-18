package config

import (
	"encoding/json"
	"fmt"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var (
	c *Config
)

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

// InitConfig 初始化配置
func InitConfig(fileName string) {

	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(fileName), "-", "_", -1))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if fileName != "" {
		viper.SetConfigFile(fileName)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("etc/")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration file(%s): %v\n", fileName, err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&c); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: Unable to decode into struct file(%s): %v\n", fileName, err)
		os.Exit(1)
	}
}

// String 配置字符输出
func (o *Config) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if c == nil {
		return &Config{}
	}
	return c
}
