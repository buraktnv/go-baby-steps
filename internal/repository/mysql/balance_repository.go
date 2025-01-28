package mysql

import (
    "context"
    "database/sql"
    "financial-service/internal/models"
    "financial-service/internal/repository"
)

type BalanceRepository struct {
    db *sql.DB
}

func NewBalanceRepository(db *sql.DB) *BalanceRepository {
    return &BalanceRepository{db: db}
}

func (r *BalanceRepository) GetBalance(ctx context.Context, userID uint) (*models.Balance, error) {
    balance := &models.Balance{}
    query := `
        SELECT user_id, amount, last_updated_at
        FROM balances WHERE user_id = ?
    `
    err := r.db.QueryRowContext(ctx, query, userID).Scan(
        &balance.UserID,
        &balance.Amount,
        &balance.LastUpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, repository.ErrNotFound
    }
    if err != nil {
        return nil, err
    }
    return balance, nil
}

func (r *BalanceRepository) UpdateBalance(ctx context.Context, balance *models.Balance) error {
    query := `
        UPDATE balances 
        SET amount = ?, last_updated_at = ?
        WHERE user_id = ?
    `
    result, err := r.db.ExecContext(ctx, query,
        balance.Amount,
        balance.LastUpdatedAt,
        balance.UserID,
    )
    if err != nil {
        return err
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rows == 0 {
        return repository.ErrNotFound
    }
    return nil
}

func (r *BalanceRepository) CreateBalance(ctx context.Context, balance *models.Balance) error {
    query := `
        INSERT INTO balances (user_id, amount, last_updated_at)
        VALUES (?, ?, ?)
    `
    _, err := r.db.ExecContext(ctx, query,
        balance.UserID,
        balance.Amount,
        balance.LastUpdatedAt,
    )
    return err
} 