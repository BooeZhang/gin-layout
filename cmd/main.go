package main

import (
	"flag"
	"fmt"
	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/core"
	"github.com/BooeZhang/gin-layout/internal/model"
	"github.com/BooeZhang/gin-layout/pkg/auth"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/store/mysqlx"
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
	st := core.NewStorageWithConfig(*cf)
	defer func() {
		st.Close()
	}()

	migrateDB(st)
	mysqlx.CreateSuperUser(st.GetMySQL(), cf.MySQLConfig)
	auth.InitAuth(cf)

	app := core.NewHttpServer(config.GetConfig())
	app.LoadRouter(initRouter(st))
	app.Run().Error()
}

func migrateDB(st *core.StoreImpl) {
	if err := st.GetMySQL().AutoMigrate(
		new(model.SysUser),
	); err != nil {
		fmt.Printf("migrate db failed: %s", err)
		os.Exit(1)
	}
	fmt.Println("migrate db completed...")
}
