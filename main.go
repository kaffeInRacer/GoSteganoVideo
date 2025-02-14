package main

import (
	"github.com/lmittmann/tint"
	"kaffein/cmd/web/server"
	"kaffein/config"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	cfg := config.Config{}
	if err := cfg.LoadEnvironment(); err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}

	app := &server.Application{
		Config: cfg,
		Logger: logger,
		Wg:     sync.WaitGroup{},
	}

	app.ServeHTTP()
}
