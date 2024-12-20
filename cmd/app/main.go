package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndreiGoStorm/previewer/internal/app"
	"github.com/AndreiGoStorm/previewer/internal/cache"
	"github.com/AndreiGoStorm/previewer/internal/config"
	"github.com/AndreiGoStorm/previewer/internal/logger"
	"github.com/AndreiGoStorm/previewer/internal/server"
	"github.com/AndreiGoStorm/previewer/internal/service"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yml", "Path to configuration file")
}

func main() {
	conf := config.New(configFile)
	logg := logger.New(conf.Log.Level)

	previewer := service.New(logg)
	lru := cache.New(conf.Cache, previewer.Storage)
	application := app.New(logg, lru, previewer, conf)

	httpServer := server.New(conf.HTTP, logg)
	httpServer.Start(application)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	select {
	case s := <-interrupt:
		logg.Info("main signal interrupt: " + s.String())
	case err := <-httpServer.Notify():
		logg.Error("main httpServer notify: %w", err)
	}

	if err := httpServer.Stop(); err != nil {
		logg.Error("main httpServer stop: %w", err)
	}
}
