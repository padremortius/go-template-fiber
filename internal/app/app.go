package app

import (
	"context"
	"go-template-fiber/internal/config"
	"go-template-fiber/internal/controller/baserouting"
	v1 "go-template-fiber/internal/controller/v1"
	"go-template-fiber/internal/crontab"
	"go-template-fiber/internal/httpserver"
	"go-template-fiber/internal/storage/sqlite"
	"go-template-fiber/internal/svclogger"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	log := svclogger.New("")

	if err := config.NewConfig(); err != nil {
		log.Logger.Fatal().Msgf("Config error: %v", err)
	}

	shutdownTimeout := config.Cfg.HTTP.Timeouts.Shutdown

	ctxParent, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Logger.Info().Msgf("Start application. Version: %v", config.Cfg.Version.Version)

	log.ChangeLogLevel(config.Cfg.Log.Level)

	// init storage
	storage, err := sqlite.New(ctxParent, config.Cfg.Storage.Path, log)
	if err != nil {
		log.Logger.Fatal().Msgf("Storage error: %v", err)
	}

	if err := storage.InitDB(); err != nil {
		log.Logger.Fatal().Msgf("Storage error: %v", err)
	}

	ctxDb := context.WithValue(ctxParent, "db", storage)

	// Init crontab
	ctb := crontab.New(ctxDb, log, &config.Cfg.Crontab)
	ctb.LoadTasks(ctxParent, &config.Cfg.Crontab)
	go ctb.StartCron()

	// HTTP Server
	log.Logger.Info().Msg("Start web-server on port " + config.Cfg.HTTP.Port)

	httpServer := httpserver.New(ctxParent, log, &config.Cfg.HTTP)
	baserouting.InitBaseRouter(httpServer.Handler)
	v1.InitAppRouter(httpServer.Handler)
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Logger.Info().Msgf("app - Run - signal: %v", s.String())
	case err := <-httpServer.Notify():
		log.Logger.Error().Msgf("app - Run - httpServer.Notify: %v", err)
	}

	// Shutdown
	if err := httpServer.Shutdown(shutdownTimeout); err != nil {
		log.Logger.Error().Msgf("app - Run - httpServer.Shutdown: %v", err)
	}
}
