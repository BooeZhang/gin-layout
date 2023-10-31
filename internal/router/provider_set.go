package router

import (
	"github.com/BooeZhang/gin-layout/server"
	"github.com/google/wire"
)

var ApiRouterProviderSet = wire.NewSet(
	NewApiRouter,
	wire.Bind(new(server.Router), new(*ApiRouter)),
)
