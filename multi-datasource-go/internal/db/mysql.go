package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver registration for database/sql
)

// OpenMySQL initializes and returns a MySQL database connection pool.
//
// Parameters:
//   - dsn:      Data Source Name containing credentials and connection info.
//     Example: "user:pass@tcp(127.0.0.1:3306)/dbname?parseTime=true"
//   - maxOpen:  Maximum number of open connections to the database.
//   - maxIdle:  Maximum number of idle (unused) connections in the pool.
//   - lifeMin:  Maximum lifetime (in minutes) for which a connection may be reused.
//   - idleMin:  Maximum idle time (in minutes) before an idle connection is closed.
//
// Returns:
//   - *sql.DB:  A configured database connection pool.
//   - error:    Non-nil if connection initialization or ping fails.
func OpenMySQL(dsn string, maxOpen, maxIdle, lifeMin, idleMin int) (*sql.DB, error) {
	// Initialize MySQL connection using provided DSN.
	db, err := sql.Open("mysql", dsn)
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
