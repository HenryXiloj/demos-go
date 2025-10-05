package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func OpenPGPool(ctx context.Context, dsn string, maxConns, minConns int, lifeMin, idleMin int) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = int32(maxConns)
	cfg.MinConns = int32(minConns)
	cfg.MaxConnLifetime = time.Duration(lifeMin) * time.Minute
	cfg.MaxConnIdleTime = time.Duration(idleMin) * time.Minute
	return pgxpool.NewWithConfig(ctx, cfg)
}
