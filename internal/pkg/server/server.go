package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/BooeZhang/gin-layout/internal/pkg/options"
	"github.com/BooeZhang/gin-layout/middleware"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// APIServer 通用的api服务.
type APIServer struct {
	Middlewares []string
	Mode        string
	RunInfo     *runInfo

	// ShutdownTimeout 优雅关闭
	ShutdownTimeout time.Duration

	*gin.Engine
	Healthz         bool
	EnableMetrics   bool
	EnableProfiling bool
	// wrapper for gin.Engine

	HttpServer *http.Server
}

// RunInfo 服务器运行配置。
type runInfo struct {
	BindAddress string
	BindPort    int
	CertKey     options.CertKey
}

// Address 将主机 IP 地址和主机端口号连接成一个地址字符串，例如：0.0.0.0:8443。
func (r *runInfo) Address() string {
	return net.JoinHostPort(r.BindAddress, strconv.Itoa(r.BindPort))
}

func InitGenericAPIServer(s *APIServer) {
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

// New 从给定的配置返回 GenericAPIServer 的新实例。
func New(opts *options.Options) *APIServer {
	s := &APIServer{
		RunInfo: &runInfo{
			BindAddress: opts.HttpServerOptions.BindAddress,
			BindPort:    opts.HttpServerOptions.BindPort,
			CertKey: options.CertKey{
				CertFile: opts.HttpServerOptions.ServerCert.CertFile,
				KeyFile:  opts.HttpServerOptions.ServerCert.KeyFile,
			},
		},
		Mode:            opts.ServerRunOptions.Mode,
		Healthz:         opts.ServerRunOptions.Healthz,
		Middlewares:     opts.ServerRunOptions.Middlewares,
		EnableMetrics:   opts.FeatureOptions.EnableMetrics,
		EnableProfiling: opts.FeatureOptions.EnableProfiling,
		Engine:          gin.New(),
	}

	InitGenericAPIServer(s)

	return s
}

// InstallAPIs 通用api。
func (s *APIServer) InstallAPIs() {
	// 添加健康检查api
	if s.Healthz {
		s.GET("/healthz", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
		})
	}

	// 添加监控
	if s.EnableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.Engine)
	}

	// 添加性能测试工具
	if s.EnableProfiling {
		pprof.Register(s.Engine)
	}

}

func (s *APIServer) Setup() {
	gin.SetMode(s.Mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

// InstallMiddlewares 安装通用中间件。
func (s *APIServer) InstallMiddlewares() {
	// necessary middlewares
	s.Use(middleware.RequestID())
	s.Use(middleware.Context())
	// s.Use(limits.RequestSizeLimiter(10))

	// install custom middlewares
	for _, m := range s.Middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)

			continue
		}

		log.Infof("install middleware: %s", m)
		s.Use(mw)
	}
}

// Run spawns the http server. It only returns when the port cannot be listened on initially.
func (s *APIServer) Run() error {
	s.HttpServer = &http.Server{
		Addr:    s.RunInfo.Address(),
		Handler: s,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}
	go func() {
		key, cert := s.RunInfo.CertKey.KeyFile, s.RunInfo.CertKey.CertFile
		if cert == "" || key == "" || s.RunInfo.BindPort == 0 {
			log.Infof("Start to listening the incoming requests on http address: %s", s.RunInfo.Address())
			if err := s.HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err.Error())
			}
			log.Infof("Server on %s stopped", s.RunInfo.Address())
		}
		log.Infof("Start to listening the incoming requests on https address: %s", s.RunInfo.Address())
		if err := s.HttpServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())
		}
		log.Infof("Server on %s stopped", s.RunInfo.Address())
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.Healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")
	if err := s.HttpServer.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: ", err)
	}
	log.Info("Server exiting")
	return nil
}

// ping 服务器健康
func (s *APIServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/healthz", s.RunInfo.Address())
	if strings.Contains(s.RunInfo.Address(), "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%d/healthz", s.RunInfo.BindPort)
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

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}
}
