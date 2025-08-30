package main

import (
	"github.com/Anacardo89/lenic/config"
	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/pkg/db"
)

func initDB(cfg config.DB) (repo.DBRepository, error) {
	pool, err := db.Connect(cfg)
	if err != nil {
		return nil, err
	}
	repo := repo.NewDBRepo(pool)
	return repo, nil
}
