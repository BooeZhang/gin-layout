package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/docs"
	middleware2 "github.com/BooeZhang/gin-layout/internal/middleware"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
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
	CertKey         config.CertKey // 是否启用https
	Middlewares     []string       // 启用的中间件
	Health          bool           // 是否启用健康检查
	EnableMetrics   bool           // 是否启用监控
	EnableProfiling bool           // 是否启用性能分析工具
	HttpServer      *http.Server
	*gin.Engine
}

// address 将主机 IP 地址和主机端口号连接成一个地址字符串，例如：0.0.0.0:8443。
func (h *HttpServer) address() string {
	return net.JoinHostPort(h.BindAddress, strconv.Itoa(h.BindPort))
}

// InitGenericAPIServer 初始化 API 服务
func InitGenericAPIServer(s *HttpServer) {
	if s.Debug {
		s.SetupSwagger()
	}
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
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

// Setup // http server 基础配置
func (h *HttpServer) Setup() {
	if !h.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
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
	if h.EnableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(h.Engine)
	}

	// 添加性能测试工具
	if h.EnableProfiling {
		pprof.Register(h.Engine)
	}

}

// InstallMiddlewares 初始化中间件。
func (h *HttpServer) InstallMiddlewares() {
	// 必要中间件
	h.Use(gin.Recovery())
	h.Use(middleware2.RequestID())
	h.Use(middleware2.Context())

	// 自定义中间件
	for _, m := range h.Middlewares {
		mw, ok := middleware2.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)

			continue
		}

		log.Infof("use middleware: %s", m)
		h.Use(mw)
	}
}

// LoadRouter 加载自定义路由
func (h *HttpServer) LoadRouter(rs ...Router) {
	for _, r := range rs {
		r.Load(h.Engine)
	}
}

// Run 启动 http 服务器.
func (h *HttpServer) Run() error {
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

	// 优雅关闭
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Info("Shutting down server...")
		if err := h.HttpServer.Shutdown(ctx); err != nil {
			log.Errorf("Server forced to shutdown: ", err)
		}
		log.Info("Server exiting")
	}()

	// 健康检查
	go func() {
		if h.Health {
			if err := h.ping(ctx); err != nil {
				log.Fatal(err.Error())
			}
		}
	}()

	key, cert := h.CertKey.KeyFile, h.CertKey.CertFile
	if cert == "" || key == "" {
		log.Infof("Start to listening the incoming requests on http address: %s", h.address())
		if err := h.HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		log.Infof("Server on %s stopped", h.address())
	} else {
		log.Infof("Start to listening the incoming requests on https address: %s", h.address())
		if err := h.HttpServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		log.Infof("Server on %s stopped", h.address())
	}

	return errors.New("service shutdown")
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
			log.Info("The router has been deployed successfully.")
			_ = resp.Body.Close()

			return nil
		}

		// 暂停 1 秒钟
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}
}

// SetupSwagger 启用swagger
func (h *HttpServer) SetupSwagger() {
	docs.SwaggerInfo.BasePath = "/v1"
	h.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
