package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrateDirection string

const (
	MigrateUp   MigrateDirection = "up"
	MigrateDown MigrateDirection = "down"
)

// migrates db with dsn and the migration path, as well as migrate direction
func MigrateDB(dsn, migrationsPath string, direction MigrateDirection) error {
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return fmt.Errorf("migration directory not found at: %s", migrationsPath)
	}
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dsn)
	if err != nil {
		return fmt.Errorf("error migrating: %s", err.Error())
	}
	if direction == MigrateUp {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	} else {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}
	return nil
}

// seeds db with file
func SeedDB(ctx context.Context, db *sql.DB, seedPath string) error {
	seed, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("failed to read seed file: %w", err)
	}
	_, err = db.Exec(string(seed))
	if err != nil {
		return fmt.Errorf("failed to execute seed: %w", err)
	}
	return nil
}
