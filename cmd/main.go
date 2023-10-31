package main

import (
	"flag"
	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/server"
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
	configFile := flag.String("c", "", "-c 选项用于指定要使用的配置文件")
	flag.Parse()

	config.InitConfig(*configFile)
	printWorkingDir()
	cf := config.GetConfig()

	log.Init(cf.LogConfig)
	mysql.InitMysql(cf.MySQLConfig)
	redis.InitRedis(cf.RedisConfig)

	app := server.NewHttpServer(config.GetConfig())
	app.LoadRouter(initRouter(mysql.GetDB(), redis.GetRedis()))
	log.Fatal(app.Run().Error())
}
