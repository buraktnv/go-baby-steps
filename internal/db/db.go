package db

import (
    "database/sql"
    "errors"
    "fmt"
    "financial-service/internal/config"
    _ "github.com/go-sql-driver/mysql"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/mysql"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewDB(cfg *config.Config) (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
        cfg.DBUser,
        cfg.DBPassword,
        cfg.DBHost,
        cfg.DBPort,
        cfg.DBName,
    )

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %w", err)
    }

    // Set connection pool settings
    db.SetMaxOpenConns(cfg.DBMaxOpenConns)
    db.SetMaxIdleConns(cfg.DBMaxIdleConns)
    db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

    // Test the connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to the database: %w", err)
    }

    return db, nil
}

func MigrateDB(db *sql.DB, cfg *config.Config) error {
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