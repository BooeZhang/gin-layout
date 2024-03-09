//go:build wireinject
// +build wireinject

package main

import (
	"github.com/BooeZhang/gin-layout/core"
	"github.com/BooeZhang/gin-layout/internal/handler"
	"github.com/BooeZhang/gin-layout/internal/repo/mysql"
	"github.com/BooeZhang/gin-layout/internal/router"
	"github.com/BooeZhang/gin-layout/internal/service"
	"github.com/BooeZhang/gin-layout/store"
	"github.com/google/wire"
)

var ApiRouterProviderSet = wire.NewSet(
	router.NewApiRouter,
	wire.Bind(new(core.Router), new(*router.ApiRouter)),
)

func initRouter(st store.Storage) core.Router {
	panic(wire.Build(
		mysql.ProviderSet,
		service.ServiceProviderSet,
		handler.ProviderSet,
		ApiRouterProviderSet,
	))
	return nil
}
