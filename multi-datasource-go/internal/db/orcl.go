package db

import (
	"database/sql"
	"time"

	_ "github.com/sijms/go-ora/v2"
)

func OpenOracle(dsn string, maxOpen, maxIdle, lifeMin, idleMin int) (*sql.DB, error) {
	db, err := sql.Open("oracle", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(time.Duration(lifeMin) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(idleMin) * time.Minute)
	return db, db.Ping()
}
