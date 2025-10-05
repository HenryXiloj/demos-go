package db

import (
	"database/sql"
	"time"

	_ "github.com/godror/godror"
)

func OpenOracle(dsn string, maxOpen, maxIdle, lifeMin, idleMin int) (*sql.DB, error) {
	db, err := sql.Open("godror", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(time.Duration(lifeMin) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(idleMin) * time.Minute)
	return db, db.Ping()
}
