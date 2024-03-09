package service

import (
	userService "github.com/BooeZhang/gin-layout/internal/service/v1/user"
	"github.com/google/wire"
)

var ServiceProviderSet = wire.NewSet(
	userService.ProviderSet,
)
