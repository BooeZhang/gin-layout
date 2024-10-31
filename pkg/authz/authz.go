package authz

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/rs/zerolog/log"

	"gin-layout/config"
	"gin-layout/store/mysqlx"
)

var (
	e *casbin.Enforcer
)

func InitAuth() {
	cf := config.GetConfig()
	a, err := gormadapter.NewAdapterByDB(mysqlx.GetDB())
	if err != nil {
		log.Fatal().Msg("初始化权限访问数据库失败")
	}

	e, err = casbin.NewEnforcer(cf.HttpServerConfig.CasbinModelPath, a)
	if err != nil {
		log.Fatal().Msg("初始化权限访问系统执行器失败")
	}
}

func GetEnforcer() *casbin.Enforcer {
	return e
}
