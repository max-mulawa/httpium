package main

import (
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

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGHUP)

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("failed to initialize logger %+v", err)
		os.Exit(1)
	}
	defer logger.Sync()
	lg := logger.Sugar()

	wd, _ := os.Getwd()
	lg.Infow("working directory", zap.String("wd", wd))

	config := config.NewHttpiumConfig()
	config.Load()

	strport := os.Getenv("HTTP_PORT")
	if strport != "" {
		val, err := strconv.ParseUint(strport, 10, 32)
		if err != nil {
			lg.Errorw("invalid port in environment variable", "invalid", strport, "err", err)
			os.Exit(4)
		}
		config.Server.Port = uint(val)
	}

	go func() {
		sig := <-sigc

		switch sig {
		case syscall.SIGTERM:
			lg.Info("sigterm received")
			os.Exit(0)
		case syscall.SIGHUP:
			lg.Info("sighup received, reloading configuration")
			//todo: reload configuration and restart server
		default:
			lg.Info("%v signal received, don't know what to do", sig)
		}

	}()

	server := server.NewServer(lg, config)
	server.Start()
}
