package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	"go.uber.org/zap"

	req "github.com/max-mulawa/httpium/pkg/http/request"
)

type HttpiumServer struct {
	port uint64
	lg   *zap.SugaredLogger
}

func NewServer(lg *zap.SugaredLogger, port uint64) *HttpiumServer {
	return &HttpiumServer{
		port: port,
		lg:   lg,
	}
}

func (s *HttpiumServer) Start() {
	lc := net.ListenConfig{
		Control:   onListeningControl,
		KeepAlive: time.Second * 10,
	}

	ln, err := lc.Listen(context.Background(), "tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.lg.Errorw("server failed to listen on port", "port", s.port)
		os.Exit(2)
	}
	defer func() {
		err := ln.Close()
		if err != nil {
			s.lg.Errorw("closing listener failed", "err", err)
		}
	}()

	host, _, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		s.lg.Errorw("failed to extract hostname from listener")
		os.Exit(3)
	}
	s.lg.Infow("server is listening for incoming connections", "host", host, "port", s.port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.lg.Errorw("Failed to accept connection", "err", err)
		}
		go handleConnection(conn, s.lg)
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

		content := string(buffer[:count])
		logger.Infow("Request received", "content", content)

		request, err := req.Parse(content)
		if err != nil {
			logger.Errorw("failed to parse request", "err", err)
			break
		}
		logger.Info(request)

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
