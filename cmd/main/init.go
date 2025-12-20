package main

import (
	"github.com/Anacardo89/lenic/config"
	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/pkg/db"
)

func initDB(cfg *config.Config) (repo.DBRepo, error) {
	if cfg.AppEnv == "aws" {
		token, err := db.GetRDSToken(cfg, cfg.DB.UserRun)
		if err != nil {
			return nil, err
		}
		cfg.DB.Pass = token
	}
	dsn, err := db.BuildDSN_URL(cfg, cfg.DB.UserRun)
	if err != nil {
		return nil, err
	}
	pool, err := db.Connect(cfg, dsn, cfg.DB.UserRun)
	if err != nil {
		return nil, err
	}
	repo := repo.NewDBRepo(pool)
	return repo, nil
}
