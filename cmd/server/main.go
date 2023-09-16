package main

import (
	"context"
	"errors"
	"log"
	"orange-backstage-api/app"
	"orange-backstage-api/infra/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

// Constants used for injecting by go build -ldflags
var (
	AppName    = "orange-backstage-api"
	AppVersion = "unknown_version"
	AppBuild   = "unknown_build"
)

var defaultCfg = config.App{
	Server: config.Server{
		RunMode:      "release",
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	},

	Log: config.Log{
		Level:        "info",
		FileName:     "server.log",
		MaxSize:      100,
		MaxBackups:   10,
		MaxAge:       30,
		Compress:     true,
		ConsoleDebug: false,
	},
}

func main() {
	viper.AddConfigPath("./")
	viper.SetConfigName("server_config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var notfoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &notfoundErr) {
			log.Println("Config file not found, use default config")
		} else {
			log.Fatalf("Error reading config file: %v", err)
		}
	}

	appCfg := defaultCfg
	if err := viper.Unmarshal(&appCfg); err != nil {
		log.Fatalf("Unable to decode into struct： %v", err)
	}

	rootCtx, cancel := context.WithCancel(context.Background())
	app, err := app.New(rootCtx, AppName, appCfg)
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
