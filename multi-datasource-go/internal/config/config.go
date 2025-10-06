package config

import (
	"github.com/spf13/viper"
)

// App holds application-level configuration parameters.
type App struct {
	// HTTPPort defines the port where the HTTP server (Gin) listens.
	// Example: 8080 or 9000
	HTTPPort int

	// RequestTimeoutSec specifies the timeout duration (in seconds)
	// for request processing to prevent long-running operations.
	RequestTimeoutSec int
}

// DB defines database configuration parameters shared by MySQL, Postgres, and Oracle.
type DB struct {
	// Enabled toggles whether this database connection should be initialized.
	Enabled bool

	// DSN (Data Source Name) contains connection info such as user, password, host, port, and database name.
	// Example:
	//   MySQL:    "user:pass@tcp(127.0.0.1:3306)/db?parseTime=true"
	//   Postgres: "postgres://user:pass@localhost:5432/db?sslmode=disable"
	//   Oracle:   "oracle://user:pass@localhost:1521/xepdb1"
	DSN string

	// MaxOpenConns defines the maximum number of open connections to the database.
	MaxOpenConns int

	// MaxIdleConns defines the maximum number of idle connections in the pool.
	MaxIdleConns int

	// ConnMaxLifetimeMin defines how long a connection may be reused before being closed (in minutes).
	ConnMaxLifetimeMin int

	// ConnMaxIdleMin defines how long an idle connection can remain before being closed (in minutes).
	ConnMaxIdleMin int
}

// Config aggregates all application and database configurations.
type Config struct {
	App      App
	MySQL    DB
	Postgres DB
	Oracle   DB
}

// Load reads configuration from application.yaml and environment variables.
// It also applies default values if certain fields are not set.
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("application") // Look for file named "application.yaml"
	v.SetConfigType("yaml")
	v.AddConfigPath(".") // Search in current directory
	v.AutomaticEnv()     // Bind environment variables automatically

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal YAML values into Config struct
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Apply default values if missing in YAML
	if cfg.App.HTTPPort == 0 {
		cfg.App.HTTPPort = 9000
	}
	if cfg.App.RequestTimeoutSec == 0 {
		cfg.App.RequestTimeoutSec = 5
	}

	return cfg, nil
}
