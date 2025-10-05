package config

import (
	"github.com/spf13/viper"
)

type App struct {
	HTTPPort          int
	RequestTimeoutSec int
}

type DB struct {
	Enabled            bool
	DSN                string
	MaxOpenConns       int
	MaxIdleConns       int
	ConnMaxLifetimeMin int
	ConnMaxIdleMin     int
}

type Config struct {
	App      App
	MySQL    DB
	Postgres DB
	Oracle   DB
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("application")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// defaults
	if cfg.App.HTTPPort == 0 {
		cfg.App.HTTPPort = 9000
	}
	if cfg.App.RequestTimeoutSec == 0 {
		cfg.App.RequestTimeoutSec = 5
	}

	return cfg, nil
}
