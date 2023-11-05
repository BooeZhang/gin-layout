package auth

import (
	"fmt"
	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var (
	e *casbin.Enforcer
)

func InitAuth(cf *config.Config) {
	uri := fmt.Sprintf("%s:%s@tcp(%s)/%s", cf.MySQLConfig.Username, cf.MySQLConfig.Password, cf.MySQLConfig.Host, cf.MySQLConfig.Database)
	a, err := gormadapter.NewAdapter("mysql", uri, true)
	if err != nil {
		log.Fatal("初始化权限访问数据库失败")
	}

	e, err = casbin.NewEnforcer(cf.CasbinConf.ModelPath, a)
	if err != nil {
		log.Fatal("初始化权限访问系统执行器失败")
	}
}

func GetEnforcer() *casbin.Enforcer {
	return e
}
