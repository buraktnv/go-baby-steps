package db

import (
    "database/sql"
    "errors"
    "fmt"
    "financial-service/internal/config"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/mysql"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB, cfg *config.Config) error {
    driver, err := mysql.WithInstance(db, &mysql.Config{})
    
    if err != nil {
        return fmt.Errorf("could not create migration driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://internal/db/migrations",
        "mysql", 
        driver,
    )

    if err != nil {
        return fmt.Errorf("could not create migrate instance: %w", err)
    }

    if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
        return fmt.Errorf("could not run migrations: %w", err)
    }

    return nil
} 