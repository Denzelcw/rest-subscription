package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"subscription/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoad()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database!")

	err = runMigrations(db, cfg.MigrationsPath)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migrations completed successfully!")

}

func runMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
