package repo

import (
	"context"

	"multi-datasource-go/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PGCompanyRepo provides the PostgreSQL-based implementation of the CompanyRepo interface.
// It encapsulates all data persistence logic for company entities using pgx connection pooling.
type PGCompanyRepo struct {
	pool *pgxpool.Pool // PostgreSQL connection pool for efficient concurrency and reuse
}

// NewPGCompanyRepo creates a new PGCompanyRepo instance using the provided pgx connection pool.
func NewPGCompanyRepo(pool *pgxpool.Pool) *PGCompanyRepo {
	return &PGCompanyRepo{pool: pool}
}

// Create inserts a new company record into the PostgreSQL 'companies' table.
// The statement uses the RETURNING clause to fetch the newly generated ID directly from the database.
// Context is used to enforce cancellation or timeout limits on the query.
// Returns the generated company ID or an error if the operation fails.
func (r *PGCompanyRepo) Create(ctx context.Context, c *domain.Company) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx,
		"INSERT INTO companies (name) VALUES ($1) RETURNING id", c.Name).
		Scan(&id)
	return id, err
}
