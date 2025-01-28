package main

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "financial-service/internal/api"
    "financial-service/internal/api/handlers"
    "financial-service/internal/config"
    "financial-service/internal/db"
    "financial-service/internal/repository/mysql"
    "financial-service/internal/services"
    
    "github.com/rs/zerolog/log"
)

func main() {
    // Load config
    cfg := config.Load()

    // Initialize DB
    database, err := db.NewDB(cfg)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to connect to database")
    }
    defer database.Close()

    // Run migrations
    if err := db.MigrateDB(database, cfg); err != nil {
        log.Fatal().Err(err).Msg("Failed to run database migrations")
    }

    // Initialize repositories
    userRepo := mysql.NewUserRepository(database)
    txRepo := mysql.NewTransactionRepository(database)
    balanceRepo := mysql.NewBalanceRepository(database)
    auditRepo := mysql.NewAuditLogRepository(database)

    // Initialize audit logger
    auditLogger := services.NewAuditLogger(auditRepo)

    // Initialize services
    userService := services.NewUserService(userRepo, balanceRepo)
    txService := services.NewTransactionService(txRepo, balanceRepo, userRepo, 5)
    balanceService := services.NewBalanceService(balanceRepo, txRepo)
    
    // Set audit loggers
    userService.SetAuditLogger(auditLogger)
    txService.SetAuditLogger(auditLogger)

    // Initialize handlers
    userHandler := handlers.NewUserHandler(userService)
    txHandler := handlers.NewTransactionHandler(txService)
    balanceHandler := handlers.NewBalanceHandler(balanceService)

    // Initialize router
    router := api.NewRouter(userHandler, txHandler, balanceHandler)

    // Create server
    srv := &http.Server{
        Addr:    ":" + cfg.ServerPort,
        Handler: router,
    }

    // Start server
    go func() {
        log.Info().Msgf("Starting server on port %s", cfg.ServerPort)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal().Err(err).Msg("Server failed")
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Info().Msg("Shutting down server...")

    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal().Err(err).Msg("Server forced to shutdown")
    }

    log.Info().Msg("Server exited properly")
}
