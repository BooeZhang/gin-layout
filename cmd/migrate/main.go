package main

import (
	"fmt"
	"github.com/BooeZhang/gin-layout/internal/apiserver/datastore/mysql"
	"github.com/BooeZhang/gin-layout/internal/apiserver/model"
	"github.com/BooeZhang/gin-layout/internal/pkg/config"
)

func main() {
	opts := config.GetConfig()
	factory, _ := mysql.GetMysqlFactoryOr(opts.MySQLConfig)
	db := factory.GetDB()
	if err := db.AutoMigrate(
		new(model.SysUserModel),
	); err != nil {
		fmt.Printf("migrate db failed: %s", err)
	}
	fmt.Println("migrate db completed...")
}
