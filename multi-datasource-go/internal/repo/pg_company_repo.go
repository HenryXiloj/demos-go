package repo

import (
	"context"

	"multi-datasource-go/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGCompanyRepo struct{ pool *pgxpool.Pool }

func NewPGCompanyRepo(pool *pgxpool.Pool) *PGCompanyRepo { return &PGCompanyRepo{pool: pool} }

func (r *PGCompanyRepo) Create(ctx context.Context, c *domain.Company) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx,
		"INSERT INTO companies (name) VALUES ($1) RETURNING id", c.Name).
		Scan(&id)
	return id, err
}
