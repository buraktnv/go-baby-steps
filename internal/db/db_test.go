package db

import (
    "database/sql"
    "testing"
    "time"
    "financial-service/internal/config"
    _ "github.com/go-sql-driver/mysql"
    "github.com/stretchr/testify/require"
)

func setupTestDatabase(t *testing.T) (*sql.DB, func()) {
    cfg := &config.Config{
        DBHost:            "127.0.0.1",
        DBPort:            "3306",
        DBUser:            "root",
        DBPassword:        "password",
        DBName:            "myproject_test",
        DBMaxOpenConns:    5,
        DBMaxIdleConns:    5,
        DBConnMaxLifetime: time.Minute,
    }

    db, err := NewDB(cfg)
    require.NoError(t, err)

    err = RunMigrations(db, cfg)
    require.NoError(t, err)

    // Initial cleanup
    cleanupTables(t, db)

    // Return cleanup function without closing the database
    return db, func() {
        cleanupTables(t, db)
        db.Close()
    }
}

// Helper function to clean tables
func cleanupTables(t *testing.T, db *sql.DB) {
    // Disable foreign key checks
    _, err := db.Exec("SET FOREIGN_KEY_CHECKS=0")
    require.NoError(t, err)

    // Clean up tables in correct order
    tables := []string{
        "transactions",
        "balances",
        "audit_logs",
        "users",
    }
    
    for _, table := range tables {
        _, err := db.Exec("TRUNCATE TABLE " + table)
        require.NoError(t, err)
    }

    // Re-enable foreign key checks
    _, err = db.Exec("SET FOREIGN_KEY_CHECKS=1")
    require.NoError(t, err)
}

func TestDatabaseConnection(t *testing.T) {
    db, cleanup := setupTestDatabase(t)
    defer cleanup()

    err := db.Ping()
    require.NoError(t, err)
}

func TestUserTableOperations(t *testing.T) {
    db, cleanup := setupTestDatabase(t)
    defer cleanup()

    // Test user insertion
    _, err := db.Exec(`
        INSERT INTO users (username, email, password_hash, role)
        VALUES ('testuser', 'test@example.com', 'hash', 'user')
    `)
    require.NoError(t, err)

    // Test duplicate email
    _, err = db.Exec(`
        INSERT INTO users (username, email, password_hash, role)
        VALUES ('testuser2', 'test@example.com', 'hash', 'user')
    `)
    require.Error(t, err, "Should not allow duplicate email")

    // Test user retrieval
    var username string
    err = db.QueryRow("SELECT username FROM users WHERE email = ?", "test@example.com").Scan(&username)
    require.NoError(t, err)
    require.Equal(t, "testuser", username)
}

func TestTransactionOperations(t *testing.T) {
    db, cleanup := setupTestDatabase(t)
    defer cleanup()

    // Create test users first
    _, err := db.Exec(`
        INSERT INTO users (username, email, password_hash, role)
        VALUES 
            ('user1', 'user1@example.com', 'hash', 'user'),
            ('user2', 'user2@example.com', 'hash', 'user')
    `)
    require.NoError(t, err)

    // Test transaction creation
    _, err = db.Exec(`
        INSERT INTO transactions (from_user_id, to_user_id, amount, type, status)
        VALUES (1, 2, 100.00, 'transfer', 'pending')
    `)
    require.NoError(t, err)

    // Test invalid user reference
    _, err = db.Exec(`
        INSERT INTO transactions (from_user_id, to_user_id, amount, type, status)
        VALUES (999, 1, 100.00, 'transfer', 'pending')
    `)
    require.Error(t, err, "Should not allow non-existent user reference")
}

func TestBalanceOperations(t *testing.T) {
    db, cleanup := setupTestDatabase(t)
    defer cleanup()

    // Create test user
    _, err := db.Exec(`
        INSERT INTO users (username, email, password_hash, role)
        VALUES ('user1', 'test_balance@example.com', 'hash', 'user')
    `)
    require.NoError(t, err)

    // Test balance creation
    _, err = db.Exec(`
        INSERT INTO balances (user_id, amount)
        VALUES (1, 100.00)
    `)
    require.NoError(t, err)

    // Test balance update
    _, err = db.Exec(`
        UPDATE balances SET amount = 200.00 WHERE user_id = 1
    `)
    require.NoError(t, err)

    // Verify balance
    var amount float64
    err = db.QueryRow("SELECT amount FROM balances WHERE user_id = 1").Scan(&amount)
    require.NoError(t, err)
    require.Equal(t, 200.00, amount)
} 