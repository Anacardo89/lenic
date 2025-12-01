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
	migratePath := filepath.Join(cfg.RootPath, "db", "migrations")
	if err := db.MigrateDB(cfg.DB.DSN, migratePath, db.MigrateUp); err != nil {
		slog.Error("failed to run up migration", "error", err)
		os.Exit(1)
	}
}
