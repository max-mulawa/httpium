package server

import (
	"context"
	"fmt"
	"net"
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
	handler HTTPHandler
	ctx     context.Context
	ln      net.Listener
}

type HTTPHandler interface {
	Handle(*req.HTTPRequest) *response.HTTPResponse
}

var (
	connectionAlive = 10 * time.Second
)

func NewServer(ctx context.Context, lg *zap.SugaredLogger, cfg *config.HttpiumConfig) *HttpiumServer {
	return &HttpiumServer{
		config:  cfg,
		lg:      lg,
		handler: static.NewStaticFiles(lg, cfg.Content.StaticDir, cfg.Content.Default),
		ctx:     ctx,
	}
}

func (s *HttpiumServer) Start() error {
	lc := net.ListenConfig{
		Control:   onListeningControl,
		KeepAlive: connectionAlive,
	}

	ln, err := lc.Listen(s.ctx, "tcp", fmt.Sprintf(":%d", s.config.Server.Port))
	if err != nil {
		return fmt.Errorf("server failed to listen on port: %d", s.config.Server.Port)
	}

	s.ln = ln

	defer func() {
		errLn := ln.Close()
		if errLn != nil {
			s.lg.Errorw("closing listener failed", "err", errLn)
		}
	}()

	host, _, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		return fmt.Errorf("failed to extract hostname from listener")
	}

	s.lg.Infow("server is listening for incoming connections", "host", host, "port", s.config.Server.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.lg.Warnw("Failed to accept connection", "err", err)

			if acceptErr, ok := err.(*net.OpError); ok && acceptErr.Op == "accept" {
				break
			}

			continue
		}

		go handleConnection(s.ctx, conn, s)
	}

	return nil
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

func handleConnection(ctx context.Context, conn net.Conn, s *HttpiumServer) {
	defer conn.Close()

	go onCancel(ctx, s.lg, conn)

	for {
		bufSize := 1024
		buffer := make([]byte, bufSize)

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

		res := s.handler.Handle(request)

		payload, err := res.Build()
		if err != nil {
			s.lg.Errorw("Failed to build response", "err", err, "response", res)
			return
		}

		count, err = conn.Write(payload)
		if err != nil {
			s.lg.Errorw("Failed to read", "err", err, "bytes", count)
			return
		}
	}
}

func onCancel(ctx context.Context, lg *zap.SugaredLogger, conn net.Conn) {
	<-ctx.Done()
	conn.Close()
	lg.Info("Closing connection on cancel")
}

func onListeningControl(network, address string, c syscall.RawConn) error {
	fmt.Println("Before listening on port, network:", network, " ,address:", address)
	return nil
}
