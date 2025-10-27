package repo

import (
	"context"
	"time"

	"github.com/Anacardo89/lenic/config"
	"github.com/Anacardo89/lenic/pkg/db"
	"github.com/Anacardo89/lenic/pkg/fs"
	"github.com/Anacardo89/lenic/pkg/testutils"
)

func InitDB(cfg config.DB) (DBRepo, error) {
	pool, err := db.Connect(cfg)
	if err != nil {
		return nil, err
	}
	repo := NewDBRepo(pool)
	return repo, nil
}

func BuildTestDBEnv(ctx context.Context) (DBRepo, func(), string, error) {
	dsn, closeDB, err := testutils.StartPostgresContainer(ctx)
	if err != nil {
		return nil, nil, "", err
	}
	cfgDB := config.DB{
		DSN:             dsn,
		MaxConns:        5,
		MinConns:        2,
		MaxConnLifetime: 30 * time.Minute,
		MaxConnIdleTime: 5 * time.Minute,
	}
	db, err := InitDB(cfgDB)
	if err != nil {
		closeDB()
		return nil, nil, "", err
	}
	migratePath, err := fs.MakeFilePath("db", "migrations")
	if err != nil {
		closeDB()
		return nil, nil, "", err
	}
	if err := testutils.MigrateDB(dsn, migratePath, testutils.MigrateUp); err != nil {
		closeDB()
		return nil, nil, "", err
	}
	seedPath, err := fs.MakeFilePath("db", "seeds")
	if err != nil {
		closeDB()
		return nil, nil, "", err
	}
	return db, closeDB, seedPath, nil
}
