package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	Pool *pgxpool.Pool
}

type PgConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string
	PoolSize int
}

func dsn(cfg PgConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)
}

func poolCoon(poolConfig *pgxpool.Config, cfg PgConfig) {
	poolConfig.MaxConns = int32(cfg.PoolSize)
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute
}

func New(ctx context.Context, cfg PgConfig) (*PostgresClient, error) {
	dsn := dsn(cfg)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgxpool config: %w", err)
	}

	poolCoon(poolConfig, cfg)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to get ping database :%w", err)
	}

	return &PostgresClient{
		Pool: pool,
	}, nil
}

func (c *PostgresClient) HealthCheck(ctx context.Context) error {
	if c.Pool == nil {
		return errors.New("database connection is not initialized")
	}

	return c.Pool.Ping(ctx)
}

func (c *PostgresClient) Close() {
	if c.Pool != nil {
		c.Pool.Close()
	}
}

func (c *PostgresClient) GetConnection(ctx context.Context) (*pgxpool.Conn, error) {
	if c.Pool == nil {
		return nil, errors.New("database conn pool is not init")
	}

	return c.Pool.Acquire(ctx)
}
