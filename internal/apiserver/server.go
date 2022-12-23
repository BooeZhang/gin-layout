package apiserver

import (
	"github.com/BooeZhang/gin-layout/internal/apiserver/router"
	"github.com/BooeZhang/gin-layout/internal/pkg/app"
	"github.com/BooeZhang/gin-layout/internal/pkg/config"
	"github.com/BooeZhang/gin-layout/internal/pkg/server"
	"github.com/BooeZhang/gin-layout/pkg/log"
)

type apiServer struct {
	redisConfig *config.RedisConfig
	mysqlConfig *config.MySQLConfig
	apiServer   *server.APIServer
}

// createAPIServer 创建api服务器
func createAPIServer(cnf *config.Config) (*apiServer, error) {
	genericServer := server.New(cnf)

	return &apiServer{
		redisConfig: cnf.RedisConfig,
		apiServer:   genericServer,
		mysqlConfig: cnf.MySQLConfig,
	}, nil
}

func (s *apiServer) Run() error {
	router.InitRouter(s.apiServer.Engine)
	return s.apiServer.Run()
}

// NewApp 创建应用程序
func NewApp() *app.App {
	cnf := config.DefaultConfig()
	application := app.NewApp(
		app.WithOptions(cnf),
		app.WithRunFunc(run(cnf)),
		app.WithName("api-server"),
	)

	return application
}

// run 创建应用程序的回调函数
func run(cnf *config.Config) app.RunFunc {
	return func() error {
		log.Init(cnf.LogConfig)
		defer log.Flush()

		fn := func(cnf *config.Config) error {
			s, err := createAPIServer(cnf)
			if err != nil {
				return err
			}

			return s.Run()
		}

		return fn(cnf)
	}
}
