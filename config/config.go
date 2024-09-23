package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/BooeZhang/gin-layout/pkg/log"
)

var (
	c *Config
)

// Config 配置选项
type Config struct {
	GRPCConfig       *GRPC        `json:"grpc"     mapstructure:"grpc"`
	HttpServerConfig *HttpServer  `json:"http"     mapstructure:"http"`
	MySQLConfig      *MySQL       `json:"mysql"    mapstructure:"mysql"`
	RedisConfig      *Redis       `json:"redis"    mapstructure:"redis"`
	MongoConfig      *Mongo       `json:"mongodb"  mapstructure:"mongodb"`
	JwtConfig        *Jwt         `json:"jwt"      mapstructure:"jwt"`
	LogConfig        *log.Options `json:"log"      mapstructure:"log"`
	CasbinConf       *Casbin      `json:"casbin"   mapstructure:"casbin"`
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
		fmt.Printf("Error: 读取配置文件失败 file(%s): %v\n", fileName, err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&c); err != nil {
		fmt.Printf("Error: 无法解码为结构体 file(%s): %v\n", fileName, err)
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
		panic("配置尚未初始化。请先调用 InitConfig 方法")
	}
	return c
}
