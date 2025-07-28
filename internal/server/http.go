package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maYkiss56/subscription-aggregation-service/internal/config"
)

type HTTPServer interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type Server struct {
	cfg     *config.Config
	httpSrv *http.Server
	handler http.Handler
}

func New(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) SetHandler(handler http.Handler) {
	s.handler = handler
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.cfg.HTTP.Host, s.cfg.HTTP.Port)
	listener, err := net.Listen(s.cfg.HTTP.Network, addr)
	if err != nil {
		log.Fatal(err)
	}

	s.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      s.handler,
		ReadTimeout:  s.cfg.HTTP.ReadTimeout,
		WriteTimeout: s.cfg.HTTP.WriteTimeout,
	}

	serveErr := make(chan error, 1)

	go func() {
		log.Printf("starting server on %s", addr)
		if err := s.httpSrv.Serve(listener); err != nil && err != http.ErrServerClosed {
			serveErr <- fmt.Errorf("server error: %w", err)
		}
		close(serveErr)
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpSrv != nil {
		return s.httpSrv.Shutdown(ctx)
	}

	return nil
}

func (s *Server) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	if err := s.Start(ctx); err != nil {
		log.Fatalf("failed to start server :%v", err)
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
