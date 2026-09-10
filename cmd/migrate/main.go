package main

import (
	"flag"
	"fmt"
	"gin-layout/internal/infra/log"
	"os"

	"gin-layout/config"
	"gin-layout/internal/infra"
	"gin-layout/internal/infra/migration"
)

func main() {
	configFile := flag.String("c", "etc/config.toml", "-c etc/config.toml")
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Printf("load config failed: %v\n", err)
		os.Exit(1)
	}

	if err := log.Init(cfg.Log); err != nil {
		panic(err)
	}
	db, err := infra.NewDatabase(&cfg.Database)
	if err != nil {
		log.Error().Err(err).Msg("connect database failed")
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("close database failed")
		}
	}()

	if err := migration.Run(db); err != nil {
		log.Error().Err(err).Msg("database migration failed")
		os.Exit(1)
	}

	log.Info().Str("driver", cfg.Database.Driver).Msg("database migration completed")
}
