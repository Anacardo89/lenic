package config

import (
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v9"
	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var configBytes []byte

func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()
	if err := yaml.Unmarshal(configBytes, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	cfg.Server.ReadTimeout *= time.Second
	cfg.Server.WriteTimeout *= time.Second
	cfg.Server.ShutdownTimeout *= time.Second
	cfg.Session.Duration *= time.Hour
	cfg.Token.Duration *= time.Minute
	cfg.DB.MaxConnLifetime *= time.Minute
	cfg.DB.MaxConnIdleTime *= time.Minute
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parsing env: %w", err)
	}
	log.Printf("final cfg: %+v", cfg)
	return cfg, nil
}

func DefaultConfig() *Config {
	return &Config{
		Server: Server{
			Port:            "8080",
			ReadTimeout:     5,  // seconds
			WriteTimeout:    10, // seconds
			ShutdownTimeout: 15, // seconds
		},
		Session: Session{
			Secret:   "session-secret",
			Duration: 24, // hours
		},
		Token: Token{
			Secret:   "token-secret",
			Duration: 60, // minutes
		},
		DB: DB{
			DSN:             "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
			MaxConns:        10,
			MinConns:        2,
			MaxConnLifetime: 30, // minutes
			MaxConnIdleTime: 5,  // minutes
		},
		Log: Log{
			Path:       "/lenic/logs",
			File:       "lenic.log",
			Level:      "info",
			MaxSize:    10, // MB
			MaxBackups: 3,
			MaxAge:     30, // days
			Compress:   true,
		},
		Mail: Mail{
			Host: "smtp.gmail.com",
			Port: 587,
			User: "example@mail.com",
			Pass: "some email app password",
		},
	}
}
