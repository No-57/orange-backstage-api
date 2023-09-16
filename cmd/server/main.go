package main

import (
	"context"
	"log"
	"orange-backstage-api/app"
	"orange-backstage-api/infra/config"
	"os"
	"os/signal"
	"syscall"
)

// Constants used for injecting by go build -ldflags
var (
	AppName    = "orange-backstage-api"
	AppVersion = "unknown_version"
	AppBuild   = "unknown_build"
)

func main() {
	rootCtx, cancel := context.WithCancel(context.Background())
	app, err := app.New(rootCtx, AppName, config.App{
		Log: config.Log{
			Level:        "debug",
			ConsoleDebug: true,
		},
	})
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	app.Run()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	cancel()
}
