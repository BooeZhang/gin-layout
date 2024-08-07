package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/viper"

	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/core"
	"github.com/BooeZhang/gin-layout/internal/model"
	"github.com/BooeZhang/gin-layout/pkg/auth"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/store/mysqlx"
)

var (
	progressMessage = color.GreenString("==>")
)

// printWorkingDir 打印工作目录
func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v 工作目录: %s", progressMessage, wd)
	log.Infof("%v 使用的配置文件为: `%s`", progressMessage, viper.ConfigFileUsed())
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
	app.Run()
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
