package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	postgresSourceData = "postgres://unn:unn@127.0.0.1:5432/unn?sslmode=disable&application_name=backend&search_path=public"
	migrationsPath     = "file://./migrations"
	driverName         = "pgx"
	databaseName       = "unn"
	connLifetime       = 10 * time.Second
)

func New() (*sqlx.DB, error) {
	conn, err := sqlx.Connect(driverName, postgresSourceData)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	conn.SetConnMaxLifetime(connLifetime)

	return conn, nil
}

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres.WithInstance: %w", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(migrationsPath, databaseName, driver)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance: %w", err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration.Up: %w", err)
	}

	return nil
}
