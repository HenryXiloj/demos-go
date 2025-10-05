package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"multi-datasource-go/internal/config"
	"multi-datasource-go/internal/db"
	"multi-datasource-go/internal/http"
	"multi-datasource-go/internal/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	// Open pools
	var (
		mysqlDB  = mustMySQL(cfg)
		pgPool   = mustPG(cfg)
		oracleDB = mustOracle(cfg)
	)

	h := &http.Handlers{
		Users:     repo.NewMySQLUserRepo(mysqlDB),
		Companies: repo.NewPGCompanyRepo(pgPool),
		Brands:    repo.NewOracleBrandRepo(oracleDB),
		Timeout:   cfg.App.ReqTimeout,
	}

	r := gin.New()
	r.Use(gin.Recovery())
	h.Register(r)

	addr := ":" + itoa(cfg.App.HTTPPort)
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func itoa(i int) string { return fmt.Sprintf("%d", i) }

func mustMySQL(cfg *config.Config) *sql.DB {
	if !cfg.MySQL.Enabled {
		return nil
	}
	dbx, err := db.OpenMySQL(cfg.MySQL.DSN, cfg.MySQL.MaxOpen, cfg.MySQL.MaxIdle, cfg.MySQL.LifeMin, cfg.MySQL.IdleMin)
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}
	return dbx
}
func mustPG(cfg *config.Config) *pgxpool.Pool {
	if !cfg.Postgres.Enabled {
		return nil
	}
	pool, err := db.OpenPGPool(context.Background(), cfg.Postgres.DSN, cfg.Postgres.MaxOpen, cfg.Postgres.MaxIdle, cfg.Postgres.LifeMin, cfg.Postgres.IdleMin)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	return pool
}
func mustOracle(cfg *config.Config) *sql.DB {
	if !cfg.Oracle.Enabled {
		return nil
	}
	dbx, err := db.OpenOracle(cfg.Oracle.DSN, cfg.Oracle.MaxOpen, cfg.Oracle.MaxIdle, cfg.Oracle.LifeMin, cfg.Oracle.IdleMin)
	if err != nil {
		log.Fatalf("oracle: %v", err)
	}
	return dbx
}
