package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Anacardo89/lenic/config"
	"github.com/Anacardo89/lenic/pkg/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	dbToken, err := db.GetRDSToken(cfg, cfg.DB.UserMigrate)
	if err != nil {
		slog.Error("failed to get RDS token", "error", err)
		os.Exit(1)
	}
	cfg.DB.Pass = dbToken
	dsn, err := db.BuildDSN_URL(cfg, cfg.DB.UserMigrate)
	if err != nil {
		slog.Error("failed to build DSN", "error", err)
		os.Exit(1)
	}
	migratePath := filepath.Join(cfg.AppHome, "migrations")
	if err := db.MigrateDB(dsn, migratePath, db.MigrateUp); err != nil {
		slog.Error("failed to run up migration", "error", err)
		os.Exit(1)
	}
	slog.Info("migrate ran successfully")
}
