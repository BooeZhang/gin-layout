package router

import (
	"github.com/BooeZhang/gin-layout/core"
	"github.com/google/wire"
)

var ApiRouterProviderSet = wire.NewSet(
	NewApiRouter,
	wire.Bind(new(core.Router), new(*ApiRouter)),
)
