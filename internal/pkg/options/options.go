package options

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
	o       *Options
	cfgFile string
)

func init() {
	pflag.StringVarP(&cfgFile, "config", "c", cfgFile, "Read configuration from specified `FILE`, "+
		"support JSON, TOML, YAML, HCL, or Java properties formats.")
}

// Options 配置选项
type Options struct {
	ServerRunOptions  *ServerRunOptions  `json:"server"   mapstructure:"server"`
	GRPCOptions       *GRPCOptions       `json:"grpc"     mapstructure:"grpc"`
	HttpServerOptions *HttpServerOptions `json:"http"     mapstructure:"http"`
	MySQLOptions      *MySQLOptions      `json:"mysql"    mapstructure:"mysql"`
	RedisOptions      *RedisOptions      `json:"redis"    mapstructure:"redis"`
	JwtOptions        *JwtOptions        `json:"jwt"      mapstructure:"jwt"`
	Log               *log.Options       `json:"log"      mapstructure:"log"`
	FeatureOptions    *FeatureOptions    `json:"feature"  mapstructure:"feature"`
}

// DefaultOption 默认配置选项
func DefaultOption() *Options {
	o = &Options{
		ServerRunOptions:  NewServerRunOptions(),
		GRPCOptions:       NewGRPCOptions(),
		HttpServerOptions: NewHttpServerOptions(),
		MySQLOptions:      NewMySQLOptions(),
		RedisOptions:      NewRedisOptions(),
		JwtOptions:        NewJwtOptions(),
		Log:               log.NewOptions(),
		FeatureOptions:    NewFeatureOptions(),
	}

	return o
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

			if names := strings.Split(basename, "-"); len(names) > 1 {
				viper.AddConfigPath("etc/")
			}

			viper.SetConfigName(basename)
		}

		if err := viper.ReadInConfig(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration file(%s): %v\n", cfgFile, err)
			os.Exit(1)
		}
		if err := viper.Unmarshal(&o); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: Unable to decode into struct file(%s): %v\n", cfgFile, err)
			os.Exit(1)
		}
	})
}

// String 配置字符输出
func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

func GetOptions() *Options {
	return o
}
