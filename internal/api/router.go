package api

import (
    "net/http"
    "financial-service/internal/api/handlers"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
    userHandler *handlers.UserHandler,
    txHandler *handlers.TransactionHandler,
    balanceHandler *handlers.BalanceHandler,
) http.Handler {
    r := chi.NewRouter()

    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // Routes
    r.Route("/api", func(r chi.Router) {
        // User routes
        r.Route("/users", func(r chi.Router) {
            r.Post("/register", userHandler.Register)
            r.Post("/login", userHandler.Login)

            // Add other user routes as needed
        })

        // Transaction routes
        r.Route("/transactions", func(r chi.Router) {
            r.Post("/credit", txHandler.Credit)
            r.Post("/debit", txHandler.Debit)
            r.Post("/transfer", txHandler.Transfer)
        })

        // Balance routes
        r.Route("/balance", func(r chi.Router) {
            r.Get("/{user_id}", balanceHandler.GetBalance)
        })
    })

    return r
} 