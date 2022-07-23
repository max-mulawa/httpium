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
	ctx     context.Context
	ln      net.Listener
}

type HttpHandler interface {
	Handle(*req.HttpRequest) *response.HttpResponse
}

func NewServer(lg *zap.SugaredLogger, config *config.HttpiumConfig, ctx context.Context) *HttpiumServer {
	return &HttpiumServer{
		config:  config,
		lg:      lg,
		handler: static.NewStaticFiles(lg, config.Content.StaticDir),
		ctx:     ctx,
	}
}

func (s *HttpiumServer) Start() {
	lc := net.ListenConfig{
		Control:   onListeningControl,
		KeepAlive: time.Second * 10,
	}

	ln, err := lc.Listen(s.ctx, "tcp", fmt.Sprintf(":%d", s.config.Server.Port))
	if err != nil {
		s.lg.Errorw("server failed to listen on port", "port", s.config.Server.Port)
		os.Exit(4)
	}
	s.ln = ln

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
			if acceptErr, ok := err.(*net.OpError); ok && acceptErr.Op == "accept" {
				break
			}
			continue
		}
		go handleConnection(conn, s.ctx, s)
	}
}

func (s *HttpiumServer) Stop() error {
	if s.ln != nil {
		if err := s.ln.Close(); err != nil {
			s.lg.Errorw("Failed to close listener", "err", err)
			return err
		}
	}
	return nil
}

func handleConnection(conn net.Conn, ctx context.Context, s *HttpiumServer) {
	defer conn.Close()
	go onCancel(s.lg, ctx, conn)
	for {
		buffer := make([]byte, 1024)
		count, err := conn.Read(buffer)
		if err != nil {
			s.lg.Errorw("Failed to read", "err", err, "bytes", count)
			return
		}

		content := string(buffer[:count])
		s.lg.Infow("Request received", "content", content)

		request, err := req.Parse(content)
		if err != nil {
			s.lg.Errorw("failed to parse request", "err", err)
			return
		}
		s.lg.Info(request)

		response := s.handler.Handle(request)

		payload, err := response.Build()
		if err != nil {
			s.lg.Errorw("Failed to build reponse", "err", err, "response", response)
			return
		}

		count, err = conn.Write([]byte(payload))
		if err != nil {
			s.lg.Errorw("Failed to read", "err", err, "bytes", count)
			return
		}
	}
}

func onCancel(lg *zap.SugaredLogger, ctx context.Context, conn net.Conn) {
	<-ctx.Done()
	conn.Close()
	lg.Info("Closing connection on cancel")
}

func onListeningControl(network, address string, c syscall.RawConn) error {
	fmt.Println("Before listening on port, network:", network, " ,address:", address)
	return nil
}
