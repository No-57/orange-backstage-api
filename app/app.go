package app

import (
	"context"
	"fmt"
	"io"
	"orange-backstage-api/app/router"
	"orange-backstage-api/app/server"
	"orange-backstage-api/app/store"
	"orange-backstage-api/app/usecase"
	"orange-backstage-api/infra/config"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"gopkg.in/natefinch/lumberjack.v2"
)

type App struct {
	Name string

	ctx context.Context
	cfg config.App
	log Logger

	store  *store.Store
	server *server.Server
}

func New(ctx context.Context, name string, cfg config.App) (*App, error) {
	app := &App{
		Name: name,
		cfg:  cfg,
	}

	app.log = app.newLogger()
	ctx = app.log.WithContext(ctx)
	app.ctx = ctx

	store, err := store.New()
	if err != nil {
		return nil, fmt.Errorf("new store: %w", err)
	}
	app.store = store

	usecase := usecase.New(app.store, usecase.Config{
		JWT: cfg.Server.JWT,
	})

	router := router.New(app.ctx, usecase, router.Param{
		Version:   "v1",
		JWT:       cfg.Server.JWT,
		EnableDoc: cfg.Server.EnableDoc,
	})

	app.server = server.New(router, app.cfg.Server)

	return app, nil
}

func (app App) newLogger() Logger {
	cfg := app.cfg.Log

	level, err := zerolog.ParseLevel(cfg.Level)
	if level == zerolog.NoLevel || err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		if cwd, err := os.Getwd(); err == nil {
			if rel, err := filepath.Rel(cwd, file); err == nil {
				file = rel
			}
		}

		return file + ":" + strconv.Itoa(line)
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano

	fileLogger := &lumberjack.Logger{
		Filename:   cfg.FileName,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	writer := []io.Writer{fileLogger}

	if cfg.ConsoleDebug {
		consoleWriter := zerolog.ConsoleWriter{
			Out: os.Stdout,

			FormatTimestamp: func(i interface{}) string {
				format := "2006-01-02T15:04:05.000000"
				now := time.Now().Local().Format(format)

				tt, ok := i.(string)
				if !ok {
					return now
				}

				t, err := time.Parse(format, tt)
				if err != nil {
					return now
				}

				return t.Local().Format(format)
			},
		}
		writer = append(writer, consoleWriter)
	}

	output := zerolog.MultiLevelWriter(writer...)
	logger := zerolog.New(output).
		With().
		Str("service", app.Name).
		Timestamp().
		Logger().
		Level(level)
	return Logger{logger}
}

func (app *App) Run() error {
	app.log.Info().Msg("running")
	defer app.log.Info().Msg("stopped")

	g := &errgroup.Group{}

	g.Go(func() error {
		log := app.log.With().
			Str("addr", app.server.Addr()).
			Str("component", "server").
			Logger()

		log.Info().Msg("running")
		if err := app.runServer(); err != nil {
			return err
		}
		log.Info().Msg("stopped")

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("running app: %w", err)
	}

	return nil
}

func (app *App) runServer() error {
	if err := app.server.Serve(); err != nil {
		return fmt.Errorf("server serve: %w", err)
	}

	return nil
}
