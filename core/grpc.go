package core

import (
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/pkg/log"
)

func NewGRPC(c *config.GRPC) *grpc.Server {
	keepParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Second * c.IdleTimeout,
		MaxConnectionAgeGrace: time.Second * c.ForceCloseWait,
		Time:                  time.Second * c.KeepAliveInterval,
		Timeout:               time.Second * c.KeepAliveTimeout,
		MaxConnectionAge:      time.Second * c.MaxLifeTime,
	})
	srv := grpc.NewServer(keepParams)

	go func() {
		lis, err := net.Listen(c.Network, c.Addr)
		if err != nil {
			log.Fatalf("GRPC failed to listen: %v", err)
		}
		log.Infof("Start grpc listen: %s", c.Addr)
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("start grpc server error: %s", err)
		}
	}()

	return srv
}
