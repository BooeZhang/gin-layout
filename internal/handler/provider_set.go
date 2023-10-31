package handler

import (
	"github.com/BooeZhang/gin-layout/internal/handler/v1/user"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	user.NewUserHandler,
)
