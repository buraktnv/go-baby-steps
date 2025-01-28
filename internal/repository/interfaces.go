package repository

import (
    "context"
    "financial-service/internal/models"
    "errors"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uint) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
}

type TransactionRepository interface {
    Create(ctx context.Context, tx *models.Transaction) error
    GetByID(ctx context.Context, id uint) (*models.Transaction, error)
    UpdateStatus(ctx context.Context, id uint, status models.TransactionStatus) error
    GetUserTransactions(ctx context.Context, userID uint, limit, offset int) ([]models.Transaction, error)
}

type BalanceRepository interface {
    GetBalance(ctx context.Context, userID uint) (*models.Balance, error)
    UpdateBalance(ctx context.Context, balance *models.Balance) error
    CreateBalance(ctx context.Context, balance *models.Balance) error
}

type AuditLogRepository interface {
    Create(ctx context.Context, log *models.AuditLog) error
    GetByEntityID(ctx context.Context, entityType string, entityID uint) ([]*models.AuditLog, error)
}

// Custom errors
var (
    ErrNotFound      = errors.New("record not found")
    ErrDuplicateKey  = errors.New("duplicate key")
    ErrInvalidData   = errors.New("invalid data")
) 