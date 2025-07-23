package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maYkiss56/subscription-aggregation-service/internal/config"
	"github.com/maYkiss56/subscription-aggregation-service/internal/server"
)

type App struct {
	cfg    *config.Config
	server *server.Server
}

func New(cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		server: server.New(cfg),
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := a.server.Start(ctx); err != nil {
			log.Printf("server error: %v", err)
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutting down due to context cancellation")
	case sig := <-sigChan:
		log.Printf("received signal: %v. Shutting down", sig)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Println("app stopped gracefully")
	return nil
}
