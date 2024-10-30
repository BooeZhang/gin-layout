package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"gin-layout/config"
	"gin-layout/core"
	"gin-layout/internal/model"
	"gin-layout/internal/router"
	"gin-layout/pkg/logger"
	"gin-layout/store/mysqlx"
)

var (
	progressMessage = "==>"
)

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Info().Msgf("%v 工作目录: %s", progressMessage, wd)
	log.Info().Msgf("%v 使用的配置文件为: `%s`", progressMessage, viper.ConfigFileUsed())
}

func main() {
	configFile := flag.String("c", "", "-c 选项用于指定要使用的配置文件")
	flag.Parse()

	config.InitConfig(*configFile)
	cf := config.GetConfig()

	logger.InitLog(cf.LogConfig.Level, cf.LogConfig.Formatter)
	printWorkingDir()

	mysqlx.DialToMysql(cf.MysqlConfig)
	defer mysqlx.Close()
	if cf.HttpServerConfig.Debug {
		migrateDB()
		mysqlx.CreateSuperUser(mysqlx.GetDB(), cf.MysqlConfig)
	}

	app := core.NewHttpServer(cf)
	app.LoadRouter(router.Admin)
	go app.Run()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-quit
		log.Info().Msgf("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if err := app.HttpServer.Shutdown(context.Background()); err != nil {
				log.Error().Err(err).Msg("Server forced to shutdown")
			}
			log.Info().Msg("Server exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func migrateDB() {
	if err := mysqlx.GetDB().AutoMigrate(
		new(model.User),
	); err != nil {
		log.Fatal().Err(err).Msg("migrate db failed")
	}
	log.Info().Msg("migrate db completed...")
}
