package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go.uber.org/zap"

	"github.com/max-mulawa/httpium/pkg/server"
)

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Resources_and_specifications
// https://datatracker.ietf.org/doc/html/rfc7230#section-2.1
func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM)

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("failed to initialize logger %+v", err)
		os.Exit(1)
	}
	defer logger.Sync()
	lg := logger.Sugar()

	port := uint64(8080)
	strport := os.Getenv("HTTP_PORT")
	if strport != "" {
		port, err = strconv.ParseUint(strport, 10, 32)
		if err != nil {
			lg.Errorw("invalid port in environment variable", "invalid", strport, "err", err)
			os.Exit(4)
		}
	}

	go func() {
		<-sigc
		lg.Info("sigterm received")
		os.Exit(0)
	}()

	server := server.NewServer(lg, port)
	server.Start()
}
