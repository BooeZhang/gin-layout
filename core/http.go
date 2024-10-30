package core

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gin-layout/config"
	"gin-layout/docs"
	"gin-layout/middleware"
	"gin-layout/pkg/logger"
)

// Router 加载路由，使用侧提供接口，实现侧需要实现该接口
type Router interface {
	Load(engine *gin.Engine)
}

// HttpServer 通用 web 服务.
type HttpServer struct {
	BindAddress     string         // 监听地址
	BindPort        int            // 监听端口
	Debug           bool           // 启动模式
	CertKey         config.CertKey // 启用https
	Middlewares     []string       // 启用的中间件
	Health          bool           // 是否启用健康检查
	EnableMetrics   bool           // 是否启用监控
	EnableProfiling bool           // 是否启用性能分析工具
	HttpServer      *http.Server   // http 服务
	*gin.Engine                    // web 驱动
}

// NewHttpServer 从给定的配置返回 GenericAPIServer 的新实例。
func NewHttpServer(cnf *config.Config) *HttpServer {
	s := &HttpServer{
		BindAddress: cnf.HttpServerConfig.BindAddress,
		BindPort:    cnf.HttpServerConfig.BindPort,
		CertKey: config.CertKey{
			CertFile: cnf.HttpServerConfig.ServerCert.CertFile,
			KeyFile:  cnf.HttpServerConfig.ServerCert.KeyFile,
		},
		Debug:           cnf.HttpServerConfig.Debug,
		Health:          cnf.HttpServerConfig.Health,
		Middlewares:     cnf.HttpServerConfig.Middlewares,
		EnableMetrics:   cnf.HttpServerConfig.EnableMetrics,
		EnableProfiling: cnf.HttpServerConfig.EnableProfiling,
		Engine:          gin.New(),
	}

	InitGenericAPIServer(s)
	return s
}

// InitGenericAPIServer 初始化 API 服务
func InitGenericAPIServer(s *HttpServer) {
	if s.Debug {
		// 启动 API 文档
		s.SetupSwagger()
	}
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

// address 将主机 IP 地址和主机端口号连接成一个地址字符串，例如：0.0.0.0:8443。
func (h *HttpServer) address() string {
	return net.JoinHostPort(h.BindAddress, strconv.Itoa(h.BindPort))
}

// Setup http server 基础配置
func (h *HttpServer) Setup() {
	if !h.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
}

// InstallMiddlewares 初始化中间件。
func (h *HttpServer) InstallMiddlewares() {
	// 必要中间件
	h.Use(gin.Recovery())
	h.Use(middleware.Cors())
	h.Use(middleware.RequestID())
	h.Use(middleware.NoCache)
	h.Use(middleware.Secure)
	h.Use(middleware.Options)
	h.Use(logger.Logger(h.Debug))

	// 自定义中间件
	for _, m := range h.Middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Warn().Msgf("can not find middleware: %s", m)
			continue
		}
		log.Info().Msgf("use middleware: %s", m)
		h.Use(mw)
	}
}

// InstallAPIs 基础 API
func (h *HttpServer) InstallAPIs() {
	// 添加健康检查api
	if h.Health {
		h.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
		})
	}

	// 添加监控
	// if h.EnableMetrics {
	// 	prometheus := ginprometheus.NewPrometheus("gin")
	// 	prometheus.Use(h.Engine)
	// }

	// 添加性能测试工具
	if h.EnableProfiling {
		pprof.Register(h.Engine)
	}

}

// LoadRouter 加载自定义路由
func (h *HttpServer) LoadRouter(rs ...Router) {
	for _, r := range rs {
		r.Load(h.Engine)
	}
}

// ping 服务器健康
func (h *HttpServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/health", h.address())
	if strings.Contains(h.address(), "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%d/health", h.BindPort)
	}

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info().Msg("The router has been deployed successfully.")
			_ = resp.Body.Close()

			return nil
		}

		// 暂停 1 秒钟
		log.Info().Msg("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal().Msg("can not ping http server within the specified time interval.")
		default:
		}
	}
}

// SetupSwagger 启用swagger
//
//go:generate swag init -g ../cmd/main.go --output ../docs
func (h *HttpServer) SetupSwagger() {
	docs.SwaggerInfo.BasePath = "/v1"
	h.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// Run 启动 http 服务器.
func (h *HttpServer) Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	h.HttpServer = &http.Server{
		Addr:           h.address(),
		Handler:        h,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 健康检查
	go func() {
		if h.Health {
			if err := h.ping(ctx); err != nil {
				log.Fatal().Err(err)
			}
		}
	}()

	log.Info().Msgf("Start to listening the incoming requests on http address: %s", h.address())

	var (
		key, cert = h.CertKey.KeyFile, h.CertKey.CertFile
		serverErr error
	)

	if cert == "" || key == "" {
		serverErr = h.HttpServer.ListenAndServe()
	} else {
		serverErr = h.HttpServer.ListenAndServeTLS(cert, key)
	}

	if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
		log.Fatal().Err(serverErr)
	}

	log.Info().Msgf("Server on %s stopped", h.address())
}
