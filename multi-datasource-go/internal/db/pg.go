package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool" // pgxpool provides a high-performance PostgreSQL connection pool
)

// OpenPGPool initializes and returns a PostgreSQL connection pool using pgxpool.
//
// Parameters:
//   - ctx:       Context used for managing cancellation and timeouts during pool creation.
//   - dsn:       Data Source Name (connection string) for PostgreSQL,
//     e.g. "postgres://user:password@localhost:5432/dbname?sslmode=disable"
//   - maxConns:  Maximum number of total connections in the pool.
//   - minConns:  Minimum number of idle connections maintained by the pool.
//   - lifeMin:   Maximum lifetime (in minutes) a connection can exist before being recycled.
//   - idleMin:   Maximum idle time (in minutes) before a connection is closed.
//
// Returns:
//   - *pgxpool.Pool: Configured PostgreSQL connection pool.
//   - error:          Non-nil if parsing the DSN or creating the pool fails.
func OpenPGPool(ctx context.Context, dsn string, maxConns, minConns int, lifeMin, idleMin int) (*pgxpool.Pool, error) {
	// Parse the provided DSN into a pgxpool configuration struct.
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// Configure connection pool parameters.
	cfg.MaxConns = int32(maxConns)                             // Set maximum allowed connections.
	cfg.MinConns = int32(minConns)                             // Maintain a minimum number of idle connections.
	cfg.MaxConnLifetime = time.Duration(lifeMin) * time.Minute // Recycle connections after this duration.
	cfg.MaxConnIdleTime = time.Duration(idleMin) * time.Minute // Close idle connections after this duration.

	// Create and return a new PostgreSQL connection pool.
	return pgxpool.NewWithConfig(ctx, cfg)
}
