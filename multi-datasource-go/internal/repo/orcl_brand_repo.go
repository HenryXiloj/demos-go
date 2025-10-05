package repo

import (
	"context"
	"database/sql"

	"multi-datasource-go/internal/domain"
)

type OracleBrandRepo struct{ db *sql.DB }

func NewOracleBrandRepo(db *sql.DB) *OracleBrandRepo { return &OracleBrandRepo{db: db} }

func (r *OracleBrandRepo) Create(ctx context.Context, b *domain.Brand) (int64, error) {
	// Oracle often uses sequences; adapt to your schema.
	// Example assumes IDENTITY column; otherwise use RETURNING INTO with godror.
	res, err := r.db.ExecContext(ctx, "INSERT INTO brands (name) VALUES (:1)", b.Name)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
