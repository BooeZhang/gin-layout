package main

import (
	"flag"
	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/internal/server"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/route"
	"github.com/BooeZhang/gin-layout/store/mysql"
	"github.com/BooeZhang/gin-layout/store/redis"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"os"
)

var (
	progressMessage = color.GreenString("==>")
)

// printWorkingDir 打印工作目录
func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
	log.Infof("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
}

func main() {
	configFile := flag.String("conf", "", "specify config file")
	flag.Parse()

	config.InitConfig(*configFile)
	cf := config.GetConfig()

	log.Init(cf.LogConfig)
	mysql.InitMysql(cf.MySQLConfig)
	redis.InitRedis(cf.RedisConfig)

	printWorkingDir()

	app := server.NewWebServer(config.GetConfig())
	route.InitRouter(app.Engine)
	log.Fatal(app.Run().Error())
}
