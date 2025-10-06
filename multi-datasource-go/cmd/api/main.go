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

// main is the application entry point.
// It loads configuration, initializes database pools, wires handlers, and starts the HTTP server.
func main() {
	// Load configuration from application.yaml and environment variables.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Open connection pools for each enabled datasource.
	var (
		mysqlDB  = mustMySQL(cfg)  // MySQL (users)
		pgPool   = mustPG(cfg)     // PostgreSQL (companies)
		oracleDB = mustOracle(cfg) // Oracle (brands)
	)

	// Create tables for local development/demo if they don't already exist.
	// This is a convenience for quick starts; remove in production.
	createTables(mysqlDB, pgPool, oracleDB)

	// Build HTTP handlers with repositories and request timeout.
	h := &http.Handlers{
		Users:     repo.NewMySQLUserRepo(mysqlDB),
		Companies: repo.NewPGCompanyRepo(pgPool),
		Brands:    repo.NewOracleBrandRepo(oracleDB),
		Timeout:   time.Duration(cfg.App.RequestTimeoutSec) * time.Second,
	}

	// Initialize Gin router and register routes.
	r := gin.New()
	// Add recovery middleware; consider adding gin.Logger() for request logs.
	r.Use(gin.Recovery())
	h.Register(r)

	// Start HTTP server on configured port.
	addr := ":" + itoa(cfg.App.HTTPPort)
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

// itoa converts an int to string using fmt.Sprintf.
// Small helper to avoid importing strconv explicitly.
func itoa(i int) string { return fmt.Sprintf("%d", i) }

// mustMySQL opens a MySQL *sql.DB pool using configuration values.
// It terminates the program if the connection cannot be established.
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

// mustPG opens a pgxpool.Pool for PostgreSQL using configuration values.
// It terminates the program if the pool cannot be created.
func mustPG(cfg *config.Config) *pgxpool.Pool {
	if !cfg.Postgres.Enabled {
		return nil
	}
	// OpenPGPool signature: (ctx, dsn, maxConns, minConns, lifeMin, idleMin)
	pool, err := db.OpenPGPool(
		context.Background(),
		cfg.Postgres.DSN,
		cfg.Postgres.MaxOpenConns,       // maxConns
		cfg.Postgres.MaxIdleConns,       // minConns (reasonable proxy from idle)
		cfg.Postgres.ConnMaxLifetimeMin, // lifeMin
		cfg.Postgres.ConnMaxIdleMin,     // idleMin
	)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	return pool
}

// mustOracle opens an Oracle *sql.DB pool using configuration values.
// It terminates the program if the connection cannot be established.
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

// createTables ensures demo/dev tables exist across all configured databases.
// For production deployments, prefer migrations managed by a tool (e.g., goose, migrate, flyway).
func createTables(mysqlDB *sql.DB, pgPool *pgxpool.Pool, oracleDB *sql.DB) {
	// Use a bounded context so DDLs don't hang indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// --- MySQL (users table) ---
	if mysqlDB != nil {
		// Create the database (no-op if it already exists); helpful for local dev.
		if _, err := mysqlDB.ExecContext(ctx, `
			CREATE DATABASE IF NOT EXISTS test_db;
		`); err != nil {
			log.Printf("mysql create db: %v", err)
		}
		// Create users table if missing.
		if _, err := mysqlDB.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS users (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				last_name VARCHAR(100) NOT NULL
			);
		`); err != nil {
			log.Printf("mysql create table: %v", err)
		}
		log.Println("✅ ensured MySQL table: users")
	}

	// --- PostgreSQL (companies table) ---
	if pgPool != nil {
		if _, err := pgPool.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS companies (
				id BIGSERIAL PRIMARY KEY,
				name TEXT NOT NULL
			);
		`); err != nil {
			log.Printf("postgres create table: %v", err)
		} else {
			log.Println("✅ ensured Postgres table: companies")
		}
	}

	// --- Oracle (brands table) ---
	if oracleDB != nil {
		// Use PL/SQL block to create table only if it doesn't exist (ignore ORA-00955).
		if _, err := oracleDB.ExecContext(ctx, `
			BEGIN
				EXECUTE IMMEDIATE 'CREATE TABLE brands (
					id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
					name VARCHAR2(100) NOT NULL
				)';
			EXCEPTION
				WHEN OTHERS THEN
					IF SQLCODE != -955 THEN RAISE; END IF; -- ORA-00955 = name is already used by an existing object
			END;
		`); err != nil {
			log.Printf("oracle create table: %v", err)
		} else {
			log.Println("✅ ensured Oracle table: brands")
		}
	}
}
