package repo

import (
	"context"

	"github.com/Anacardo89/lenic/pkg/fs"
	"github.com/Anacardo89/lenic/pkg/testutils"
)

func InitDB(ctx context.Context, dsn string) (DBRepo, error) {
	pool, err := testutils.ConnectDB(ctx, dsn)
	if err != nil {
		return nil, err
	}
	repo := NewDBRepo(pool)
	return repo, nil
}

func BuildTestDBEnv(ctx context.Context) (DBRepo, string, func(), string, error) {
	dsn, closeDB, err := testutils.StartPostgresContainer(ctx)
	if err != nil {
		return nil, "", nil, "", err
	}
	db, err := InitDB(ctx, dsn)
	if err != nil {
		closeDB()
		return nil, "", nil, "", err
	}
	migratePath, err := fs.MakeFilePath("db", "migrations")
	if err != nil {
		closeDB()
		return nil, "", nil, "", err
	}
	if err := testutils.MigrateDB(dsn, migratePath, testutils.MigrateUp); err != nil {
		closeDB()
		return nil, "", nil, "", err
	}
	seedPath, err := fs.MakeFilePath("db", "seeds")
	if err != nil {
		closeDB()
		return nil, "", nil, "", err
	}
	return db, dsn, closeDB, seedPath, nil
}

func SeedDB(ctx context.Context, dsn, seed string) error {
	db, err := testutils.ConnectDB(ctx, dsn)
	if err != nil {
		return err
	}
	if err := testutils.SeedDB(ctx, db, seed); err != nil {
		return err
	}
	return nil
}
