package repo

import (
	"context"
	"database/sql"

	"multi-datasource-go/internal/domain"
)

// OracleBrandRepo provides the Oracle-based implementation of the BrandRepo interface.
// It handles all brand-related persistence operations using an Oracle database.
type OracleBrandRepo struct {
	db *sql.DB // Shared connection pool for Oracle database connections
}

// NewOracleBrandRepo creates and returns a new instance of OracleBrandRepo.
// The caller provides a *sql.DB connection already configured for Oracle.
func NewOracleBrandRepo(db *sql.DB) *OracleBrandRepo {
	return &OracleBrandRepo{db: db}
}

// Create inserts a new brand record into the Oracle 'brands' table.
// The SQL statement uses Oracle-style positional bind parameters (:1).
// Depending on your Oracle schema, you might replace this with:
//
//	"INSERT INTO brands (id, name) VALUES (brand_seq.NEXTVAL, :1)"
//
// or use RETURNING INTO if you need to fetch the generated ID.
// Returns the number of affected rows or an error if the operation fails.
func (r *OracleBrandRepo) Create(ctx context.Context, b *domain.Brand) (int64, error) {
	// Execute the INSERT command within the provided context (supports timeout/cancel)
	res, err := r.db.ExecContext(ctx, "INSERT INTO brands (name) VALUES (:1)", b.Name)
	if err != nil {
		return 0, err
	}

	// Get number of rows affected (should be 1 if successful)
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
