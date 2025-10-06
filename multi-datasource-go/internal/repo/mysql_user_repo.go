package repo

import (
	"context"
	"database/sql"
	"multi-datasource-go/internal/domain"
)

// MySQLUserRepo provides the MySQL-based implementation of the UserRepo interface.
// It encapsulates all database operations for managing User records.
type MySQLUserRepo struct {
	db *sql.DB // Shared connection pool to the MySQL database
}

// NewMySQLUserRepo creates a new MySQLUserRepo instance.
// The caller provides a configured *sql.DB connection pool.
func NewMySQLUserRepo(db *sql.DB) *MySQLUserRepo {
	return &MySQLUserRepo{db: db}
}

// Create inserts a new user record into the MySQL 'users' table.
// It accepts a context for cancellation and timeout control.
// Returns the newly inserted record ID, or an error if the insert fails.
func (r *MySQLUserRepo) Create(ctx context.Context, u *domain.User) (int64, error) {
	// Execute the INSERT statement using a prepared query with parameter placeholders (safe from SQL injection)
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO users (name, last_name) VALUES (?, ?)", u.Name, u.LastName)
	if err != nil {
		return 0, err
	}

	// Retrieve the last inserted ID (auto-increment primary key)
	return res.LastInsertId()
}
