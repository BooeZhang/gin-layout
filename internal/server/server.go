package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-layout/internal/infra"
	"gin-layout/internal/middleware"
)

// Config 保存 HTTP 服务器配置
type Config struct {
	Host string
	Port int
	Mode string
}

// Router 定义了路由注册的接口
type Router interface {
	Register(*gin.Engine)
}

// Server 封装了 Gin 引擎和 HTTP 服务器
type Server struct {
	engine *gin.Engine
	server *http.Server
	logger *infra.Logger
}

// NewServer 创建一个带有中间件和路由的 Server
func NewServer(cfg Config, logger *infra.Logger, routers ...Router) *Server {
	gin.SetMode(cfg.Mode)
	engine := gin.New()
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Recovery(logger))
	engine.Use(middleware.ErrorHandler(logger))
	engine.Use(middleware.RequestLogger(logger))
	engine.Use(middleware.CORS())

	for _, r := range routers {
		r.Register(engine)
	}

	return &Server{
		engine: engine,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler: engine,
		},
		logger: logger,
	}
}

// Start 启动 HTTP server
func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown 优雅关闭
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
