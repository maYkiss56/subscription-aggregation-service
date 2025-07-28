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
	"github.com/maYkiss56/subscription-aggregation-service/internal/delivery/api/sub"
	"github.com/maYkiss56/subscription-aggregation-service/internal/repository"
	"github.com/maYkiss56/subscription-aggregation-service/internal/server"
	"github.com/maYkiss56/subscription-aggregation-service/internal/service"
	"github.com/maYkiss56/subscription-aggregation-service/pkg/client/postgresql"
)

type App struct {
	cfg      *config.Config
	server   *server.Server
	pgClient *postgresql.PostgresClient
}

func New(cfg *config.Config) (*App, error) {
	pgCfg := postgresql.PgConfig{
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Database: cfg.Postgres.Database,
		SSLMode:  cfg.Postgres.SSLMode,
		PoolSize: cfg.Postgres.PoolSize,
	}

	pgClient, err := postgresql.New(context.Background(), pgCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres client: %w", err)
	}

	subRepo := repository.New(pgClient)

	subService := service.New(subRepo)

	subHandler := sub.New(subService)

	router := sub.NewRouter(subHandler)

	srv := server.New(cfg)
	srv.SetHandler(router)

	return &App{
		cfg:      cfg,
		server:   srv,
		pgClient: pgClient,
	}, nil
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go a.databaseHealthCheck(ctx)

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

	a.pgClient.Close()

	log.Println("app stopped gracefully")
	return nil
}

func (a *App) databaseHealthCheck(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := a.pgClient.HealthCheck(ctx); err != nil {
				log.Printf("database health check failed: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
