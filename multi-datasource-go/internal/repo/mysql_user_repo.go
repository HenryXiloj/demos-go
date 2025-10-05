package repo

import (
	"context"
	"database/sql"
	"multi-datasource-go/internal/domain"
)

type MySQLUserRepo struct{ db *sql.DB }

func NewMySQLUserRepo(db *sql.DB) *MySQLUserRepo { return &MySQLUserRepo{db: db} }

func (r *MySQLUserRepo) Create(ctx context.Context, u *domain.User) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO users (name, last_name) VALUES (?, ?)", u.Name, u.LastName)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
