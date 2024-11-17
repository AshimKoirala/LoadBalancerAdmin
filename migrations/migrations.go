package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies pending migrations.
func RunMigrations(sqldb *sql.DB) error {
	workingDir, _ := os.Getwd()
	fmt.Print(workingDir)
	migrationsPath := filepath.Join(workingDir, "./migrations")

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) { //ensuring path
		return fmt.Errorf("migrations folder not found at %s", migrationsPath)
	}

	sourcePath := fmt.Sprintf("file://%s", migrationsPath)

	// Create the Postgres driver
	driver, err := postgres.WithInstance(sqldb, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create Postgres driver: %w", err)
	}

	// Create the migration instance
	m, err := migrate.NewWithDatabaseInstance(
		sourcePath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("Looking for migrations at:", migrationsPath)

	fmt.Println("Migrations applied successfully.")
	return nil
}
