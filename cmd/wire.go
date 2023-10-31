//go:build wireinject
// +build wireinject

package main

import (
	"github.com/BooeZhang/gin-layout/internal/handler"
	"github.com/BooeZhang/gin-layout/internal/repo/mysql"
	"github.com/BooeZhang/gin-layout/internal/router"
	userService "github.com/BooeZhang/gin-layout/internal/service/v1/user"
	"github.com/BooeZhang/gin-layout/server"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var ApiRouterProviderSet = wire.NewSet(
	router.NewApiRouter,
	wire.Bind(new(server.Router), new(*router.ApiRouter)),
)

func initRouter(ds *gorm.DB, rs redis.UniversalClient) server.Router {
	panic(wire.Build(
		mysql.ProviderSet,
		userService.ProviderSet,
		handler.ProviderSet,
		ApiRouterProviderSet,
	))
	return nil
}
