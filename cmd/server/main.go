package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-layout/config"
	"gin-layout/internal/bootstrap"
)

func main() {
	configFile := flag.String("c", "etc/config.toml", "-c etc/config.toml")
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Printf("load config failed: %v\n", err)
		os.Exit(1)
	}

	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		fmt.Printf("wire failed: %v\n", err)
		os.Exit(1)
	}
	defer app.Cleanup()

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- app.HTTPServer.Start()
	}()

	app.Logger.Info().
		Str("host", cfg.Server.Host).
		Int("port", cfg.Server.Port).
		Msg("server started, waiting for shutdown signal")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case err := <-serverErr:
		if err != nil {
			app.Logger.Error().Err(err).Msg("server start failed")
		}
		return
	case sig := <-quit:
		app.Logger.Info().Str("signal", sig.String()).Msg("shutdown signal received")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.HTTPServer.Shutdown(ctx); err != nil {
		app.Logger.Error().Err(err).Msg("server shutdown failed")
	}

	if err := <-serverErr; err != nil {
		app.Logger.Error().Err(err).Msg("server stopped with error")
	}
}
