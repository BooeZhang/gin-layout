package main

import (
	"fmt"
	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/internal/model"
	"github.com/BooeZhang/gin-layout/store/mysql"
)

func main() {
	opts := config.GetConfig()
	mysql.InitMysql(opts.MySQLConfig)
	db := mysql.GetDB()
	if err := db.AutoMigrate(
		new(model.User),
	); err != nil {
		fmt.Printf("migrate db failed: %s", err)
	}
	fmt.Println("migrate db completed...")
}
