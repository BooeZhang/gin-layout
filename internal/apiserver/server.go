package apiserver

import (
	"github.com/BooeZhang/gin-layout/internal/apiserver/router"
	"github.com/BooeZhang/gin-layout/internal/pkg/app"
	"github.com/BooeZhang/gin-layout/internal/pkg/options"
	"github.com/BooeZhang/gin-layout/internal/pkg/server"
	"github.com/BooeZhang/gin-layout/pkg/log"
)

type apiServer struct {
	redisOptions *options.RedisOptions
	mysqlOptions *options.MySQLOptions
	// etcdOptions      *genericoptions.EtcdOptions
	gRPCAPIServer *grpcAPIServer
	apiServer     *server.APIServer
}

// createAPIServer 创建api服务器
func createAPIServer(opts *options.Options) (*apiServer, error) {
	genericServer := server.New(opts)

	return &apiServer{
		redisOptions: opts.RedisOptions,
		apiServer:    genericServer,
		mysqlOptions: opts.MySQLOptions,
		// etcdOptions:      cfg.EtcdOptions,
	}, nil
}

func (s *apiServer) Run() error {
	router.InitRouter(s.apiServer.Engine)
	return s.apiServer.Run()
}

// NewApp 创建应用程序
func NewApp() *app.App {
	opts := options.DefaultOption()
	application := app.NewApp(
		app.WithOptions(opts),
		app.WithRunFunc(run(opts)),
	)

	return application
}

// run 创建应用程序的回调函数
func run(opts *options.Options) app.RunFunc {
	return func() error {
		log.Init(opts.Log)
		defer log.Flush()

		fn := func(opts *options.Options) error {
			server, err := createAPIServer(opts)
			if err != nil {
				return err
			}

			return server.Run()
		}

		return fn(opts)
	}
}
