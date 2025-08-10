package db

import (
	"context"
	"fmt"
	"time"

	"github.com/Anacardo89/lenic/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbClient struct {
	cfg  *config.DBConfig
	pool *pgxpool.Pool
}

func Connect(ctx context.Context, cfg *config.DBConfig) (DBRepository, error) {
	dsn := fmt.Sprintf(cfg.DSN, cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)

	dbCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	dbCfg.MaxConns = 10
	dbCfg.MinConns = 2
	dbCfg.MaxConnLifetime = time.Hour
	dbCfg.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	dbClient := &dbClient{
		pool: pool,
	}

	return dbClient, nil
}
