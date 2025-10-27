package testutils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	testcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type MigrateDirection string

const (
	MigrateUp   MigrateDirection = "up"
	MigrateDown MigrateDirection = "down"
)

// starts a Postgres container and returns the dsn + close function.
func StartPostgresContainer(ctx context.Context) (string, func(), error) {
	// run container
	pgContainer, err := testcontainer.Run(ctx, "postgres:16",
		testcontainer.WithDatabase("testdb"),
		testcontainer.WithUsername("test"),
		testcontainer.WithPassword("secret"),
	)
	if err != nil {
		return "", nil, err
	}
	// get dsn
	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		pgContainer.Terminate(ctx)
		return "", nil, err
	}
	// close function
	close := func() { _ = pgContainer.Terminate(ctx) }
	return dsn, close, nil
}

// connects to db with dsn and returns a pool
func ConnectDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	return pool, nil
}

// migrates db with dsn and the migration path, as well as migrate direction
func MigrateDB(dsn, migrationsPath string, direction MigrateDirection) error {

	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if direction == MigrateUp {
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No migrations to apply.")
				return nil
			}
			return fmt.Errorf("migration failed: %w", err)
		}
	} else {
		if err := m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No migrations to apply.")
				return nil
			}
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	log.Println("Database migrated successfully.")
	return nil
}

// seeds db with file
func SeedDB(ctx context.Context, pool *pgxpool.Pool, seedPath string) error {
	// parse file
	seedSQL, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("failed to read seed file: %w", err)
	}
	// seed db
	_, err = pool.Exec(ctx, string(seedSQL))
	if err != nil {
		return fmt.Errorf("failed to execute seed SQL: %w", err)
	}
	return nil
}
