package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"multi-datasource-go/internal/config"
	"multi-datasource-go/internal/db"
	"multi-datasource-go/internal/http"
	"multi-datasource-go/internal/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

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
		Timeout:   time.Duration(cfg.App.RequestTimeoutSec) * time.Second,
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
	dbx, err := db.OpenMySQL(
		cfg.MySQL.DSN,
		cfg.MySQL.MaxOpenConns,
		cfg.MySQL.MaxIdleConns,
		cfg.MySQL.ConnMaxLifetimeMin,
		cfg.MySQL.ConnMaxIdleMin,
	)
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}
	return dbx
}

func mustPG(cfg *config.Config) *pgxpool.Pool {
	if !cfg.Postgres.Enabled {
		return nil
	}
	// OpenPGPool signature: (ctx, dsn, maxConns, minConns, lifeMin, idleMin)
	pool, err := db.OpenPGPool(
		context.Background(),
		cfg.Postgres.DSN,
		cfg.Postgres.MaxOpenConns,       // maxConns
		cfg.Postgres.MaxIdleConns,       // use idle as a reasonable minConns
		cfg.Postgres.ConnMaxLifetimeMin, // lifeMin
		cfg.Postgres.ConnMaxIdleMin,     // idleMin
	)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	return pool
}

func mustOracle(cfg *config.Config) *sql.DB {
	if !cfg.Oracle.Enabled {
		return nil
	}
	dbx, err := db.OpenOracle(
		cfg.Oracle.DSN,
		cfg.Oracle.MaxOpenConns,
		cfg.Oracle.MaxIdleConns,
		cfg.Oracle.ConnMaxLifetimeMin,
		cfg.Oracle.ConnMaxIdleMin,
	)
	if err != nil {
		log.Fatalf("oracle: %v", err)
	}
	return dbx
}
