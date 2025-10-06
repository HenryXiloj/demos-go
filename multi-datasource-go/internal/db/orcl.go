package db

import (
	"database/sql"
	"time"

	_ "github.com/sijms/go-ora/v2" // Oracle driver registration for database/sql
)

// OpenOracle initializes and returns an Oracle database connection pool.
//
// Parameters:
//   - dsn:      Data Source Name (connection string) that includes credentials and host information.
//     Example: "oracle://user:password@127.0.0.1:1521/xepdb1"
//   - maxOpen:  Maximum number of open connections to the database.
//   - maxIdle:  Maximum number of idle (unused) connections retained in the pool.
//   - lifeMin:  Maximum lifetime (in minutes) that a single connection can remain open.
//   - idleMin:  Maximum idle time (in minutes) before an idle connection is closed.
//
// Returns:
//   - *sql.DB:  A configured Oracle database connection pool.
//   - error:    Non-nil if the connection initialization or ping test fails.
func OpenOracle(dsn string, maxOpen, maxIdle, lifeMin, idleMin int) (*sql.DB, error) {
	// Initialize Oracle connection using provided DSN.
	db, err := sql.Open("oracle", dsn)
	if err != nil {
		return nil, err
	}

	// Configure connection pool settings.
	db.SetMaxOpenConns(maxOpen)                                 // Limit total number of open connections.
	db.SetMaxIdleConns(maxIdle)                                 // Limit number of idle connections retained.
	db.SetConnMaxLifetime(time.Duration(lifeMin) * time.Minute) // Set maximum lifetime for a single connection.
	db.SetConnMaxIdleTime(time.Duration(idleMin) * time.Minute) // Set maximum idle time before closing a connection.

	// Verify the database connection.
	return db, db.Ping()
}
