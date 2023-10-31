package user

import (
	v1 "github.com/BooeZhang/gin-layout/internal/service/v1"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	v1.NewServiceContext,
	NewUserService,
	wire.Bind(new(Service), new(*serviceImpl)),
)
