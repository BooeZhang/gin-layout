package apiserver

import (
	"github.com/BooeZhang/gin-layout/internal/apiserver/router"
	"github.com/BooeZhang/gin-layout/internal/pkg/app"
	"github.com/BooeZhang/gin-layout/internal/pkg/config"
	"github.com/BooeZhang/gin-layout/internal/pkg/server"
	"github.com/BooeZhang/gin-layout/pkg/log"
)

type mainServer struct {
	redisConfig *config.RedisConfig
	mysqlConfig *config.MySQLConfig
	apiServer   *server.APIServer
}

// NewApp 创建应用程序
func NewApp(baseName string) *app.App {
	cnf := config.DefaultConfig()
	application := app.NewApp(
		app.WithOptions(cnf),
		app.WithRunFunc(run(cnf)),
		app.WithName(baseName),
	)

	return application
}

// run 创建应用程序的回调函数
func run(cnf *config.Config) app.RunFunc {
	return func() error {
		log.Init(cnf.LogConfig)
		defer log.Flush()

		fn := func(cnf *config.Config) error {
			s, err := createMainServer(cnf)
			if err != nil {
				return err
			}

			return s.Run()
		}

		return fn(cnf)
	}
}

// createMainServer 创建主服务
func createMainServer(cnf *config.Config) (*mainServer, error) {
	genericServer := server.New(cnf)

	return &mainServer{
		redisConfig: cnf.RedisConfig,
		apiServer:   genericServer,
		mysqlConfig: cnf.MySQLConfig,
	}, nil
}

func (s *mainServer) Run() error {
	router.InitRouter(s.apiServer.Engine)
	return s.apiServer.Run()
}
