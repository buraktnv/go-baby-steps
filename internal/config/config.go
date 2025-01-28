package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    // Database configuration
    DBUser           string
    DBPassword       string
    DBHost           string
    DBPort           string
    DBName           string
    DBMaxOpenConns   int
    DBMaxIdleConns   int
    DBConnMaxLifetime time.Duration

    // Server configuration
    ServerPort string
}

func Load() *Config {
    return &Config{
        // Database configuration
        DBUser:     getEnv("DB_USER", "root"),
        DBPassword: getEnv("DB_PASSWORD", "password"),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "3306"),
        DBName:     getEnv("DB_NAME", "financial_service"),
        
        // Database pool configuration
        DBMaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
        DBMaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
        DBConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),

        // Server configuration
        ServerPort: getEnv("SERVER_PORT", "8080"),
    }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := os.Getenv(key)
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    valueStr := os.Getenv(key)
    if value, err := time.ParseDuration(valueStr); err == nil {
        return value
    }
    return defaultValue
} 