package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Resources_and_specifications
func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("failed to initialize logger %+v", err)
		os.Exit(1)
	}
	defer logger.Sync()
	suggar := logger.Sugar()

	port := 8080

	lc := net.ListenConfig{
		Control:   onListeningControl,
		KeepAlive: time.Second * 10,
	}

	ln, err := lc.Listen(context.Background(), "tcp", fmt.Sprintf(":%d", port))
	//ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		suggar.Errorw("server failed to listen on port", "port", port)
		os.Exit(2)
	}
	defer func() {
		err := ln.Close()
		if err != nil {
			suggar.Errorw("closing listener failed", "err", err)
		}
	}()

	host, _, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		suggar.Errorw("failed to extract hostname from listener")
		os.Exit(3)
	}
	suggar.Infow("server is listening for incoming connections", "host", host, "port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			suggar.Errorw("Failed to accept connection", "err", err)
		}
		go handleConnection(conn, suggar)
	}
}

func handleConnection(conn net.Conn, logger *zap.SugaredLogger) {
	for {
		buffer := make([]byte, 1024)
		count, err := conn.Read(buffer)
		if err != nil {
			logger.Errorw("Failed to read", "err", err, "bytes", count)
			break
		}

		request := string(buffer[:count])
		logger.Infow("Request received", "content", request)

		requestLines := strings.Split(request, "\r\n")
		firstLine := requestLines[0]

		logger.Info(firstLine)

		count, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 12\r\nContent-Type: text/plain\r\n\r\nHello World!"))
		if err != nil {
			logger.Errorw("Failed to read", "err", err, "bytes", count)
			break
		}
	}
	conn.Close()
}

func onListeningControl(network, address string, c syscall.RawConn) error {
	fmt.Println("Before listening on port, network:", network, " ,address:", address)
	return nil
}
