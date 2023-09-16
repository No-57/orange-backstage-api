package app

import (
	"context"
	"io"
	"orange-backstage-api/infra/config"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type App struct {
	Name string

	ctx context.Context
	cfg config.App
	log Logger
}

func New(ctx context.Context, name string, cfg config.App) (*App, error) {
	app := &App{
		Name: name,
		cfg:  cfg,
	}

	app.log = app.newLogger()
	ctx = app.log.WithContext(ctx)
	app.ctx = ctx

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

func (app *App) Run() {
	app.log.Info().Msg("running")
	defer app.log.Info().Msg("stopped")

	// TODO: implement app running logic
}
