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
	"github.com/max-mulawa/httpium/pkg/http/response"
	"github.com/max-mulawa/httpium/pkg/server/config"
	"github.com/max-mulawa/httpium/pkg/server/static"
)

type HttpiumServer struct {
	config  *config.HttpiumConfig
	lg      *zap.SugaredLogger
	handler HttpHandler
}

type HttpHandler interface {
	Handle(*req.HttpRequest) *response.HttpResponse
}

func NewServer(lg *zap.SugaredLogger, config *config.HttpiumConfig) *HttpiumServer {
	return &HttpiumServer{
		config:  config,
		lg:      lg,
		handler: static.NewStaticFiles(lg, config.Content.StaticDir),
	}
}

func (s *HttpiumServer) Start() {
	lc := net.ListenConfig{
		Control:   onListeningControl,
		KeepAlive: time.Second * 10,
	}

	ln, err := lc.Listen(context.Background(), "tcp", fmt.Sprintf(":%d", s.config.Server.Port))
	if err != nil {
		s.lg.Errorw("server failed to listen on port", "port", s.config.Server.Port)
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
	s.lg.Infow("server is listening for incoming connections", "host", host, "port", s.config.Server.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.lg.Errorw("Failed to accept connection", "err", err)
		}
		go handleConnection(conn, s)
	}
}

func handleConnection(conn net.Conn, s *HttpiumServer) {
	for {
		buffer := make([]byte, 1024)
		count, err := conn.Read(buffer)
		if err != nil {
			s.lg.Errorw("Failed to read", "err", err, "bytes", count)
			break
		}

		content := string(buffer[:count])
		s.lg.Infow("Request received", "content", content)

		request, err := req.Parse(content)
		if err != nil {
			s.lg.Errorw("failed to parse request", "err", err)
			break
		}
		s.lg.Info(request)

		response := s.handler.Handle(request)

		payload, err := response.Build()
		if err != nil {
			s.lg.Errorw("Failed to build reponse", "err", err, "response", response)
			break
		}

		count, err = conn.Write([]byte(payload))
		if err != nil {
			s.lg.Errorw("Failed to read", "err", err, "bytes", count)
			break
		}
	}
	conn.Close()
}

func onListeningControl(network, address string, c syscall.RawConn) error {
	fmt.Println("Before listening on port, network:", network, " ,address:", address)
	return nil
}
