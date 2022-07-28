package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go.uber.org/zap"

	"github.com/max-mulawa/httpium/pkg/server"
	"github.com/max-mulawa/httpium/pkg/server/config"
)

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Resources_and_specifications
// https://datatracker.ietf.org/doc/html/rfc7230#section-2.1
// https://developer.mozilla.org/en-US/docs/Web/HTTP

var (
	success                   = 0
	failedToStop              = 1
	failedToInitializeLogging = 2
)

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("failed to initialize logger %+v", err)
		os.Exit(failedToInitializeLogging)
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("faling to sync logger: %v", err)
		}
	}()

	lg := logger.Sugar()
	wd, _ := os.Getwd()
	lg.Infow("working directory", zap.String("wd", wd))

	for {
		ctx, cancelCtx := context.WithCancel(context.Background())
		cfg := loadConfiguration(lg)
		srv := server.NewServer(ctx, lg, cfg)

		go onSignal(lg, srv, sigc, cancelCtx)

		if err := srv.Start(); err != nil {
			lg.Fatalw("failed to start", "err", err)
		}
	}
}

func loadConfiguration(lg *zap.SugaredLogger) *config.HttpiumConfig {
	cfg := config.NewHttpiumConfig()
	if err := cfg.Load(); err != nil {
		lg.Fatalw("failed to load configuration", "err", err)
	}

	strport := os.Getenv("HTTP_PORT")
	if strport != "" {
		base := 10
		size := 32

		val, err := strconv.ParseUint(strport, base, size)
		if err != nil {
			lg.Fatalw("invalid port in environment variable", "invalid", strport, "err", err)
		}

		cfg.Server.Port = uint(val)
	}

	return cfg
}

func onSignal(lg *zap.SugaredLogger, srv *server.HttpiumServer, sigc chan os.Signal, cancelCtx context.CancelFunc) {
	sig := <-sigc

	switch sig {
	case syscall.SIGTERM:
	case syscall.SIGINT:
		lg.Info("sigterm received")
		stop(srv, cancelCtx)
		os.Exit(success)
	case syscall.SIGHUP:
		lg.Info("sighup received, reloading configuration")
		stop(srv, cancelCtx)
	default:
		lg.Info("%v signal received, don't know what to do", sig)
	}
}

func stop(srv *server.HttpiumServer, cancelCtx context.CancelFunc) {
	err := srv.Stop()

	cancelCtx()

	if err != nil {
		os.Exit(failedToStop)
	}
}
