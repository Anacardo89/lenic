package config

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Anacardo89/lenic/pkg/testutils"
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()
	appEnv := os.Getenv("APP_ENV")
	var rootPath string
	if appEnv == "local" {
		var err error
		rootPath, err = testutils.GetProjectRoot()
		if err != nil {
			log.Println("failed to get project root")
			return nil, err
		}
	}
	cfg.RootPath = rootPath
	if err := godotenv.Load(filepath.Join(rootPath, ".env")); err != nil {
		log.Println("No env file found, relying on default env variables")
	}
	cfgPath := os.Getenv("CFG_PATH")
	cfgFile := os.Getenv("CFG_FILE")
	cfgPath = filepath.Join(rootPath, cfgPath, cfgFile)
	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
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
			Host:            "localhost",
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
			DSN:             "postgres://user:pass@db:5432/dbname?sslmode=disable",
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
		Img: Img{
			BasePath:      "/lenic/img",
			OriginalsDir:  "/originals",
			PreviewsDir:   "/previews",
			PreviewWidth:  200, // pixels
			PreviewHeight: 200, // pixels
			JPEGQuality:   95,
		},
		Mail: Mail{
			Host: "smtp.gmail.com",
			Port: 587,
			User: "example@mail.com",
			Pass: "some email app password",
		},
	}
}
