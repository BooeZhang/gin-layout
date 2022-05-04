package apiserver

import (
	"github.com/BooeZhang/gin-layout/internal/apiserver/router"
	"github.com/BooeZhang/gin-layout/internal/pkg/options"
	"github.com/BooeZhang/gin-layout/internal/pkg/server"
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
