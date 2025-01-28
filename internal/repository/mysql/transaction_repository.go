package mysql

import (
    "context"
    "database/sql"
    "financial-service/internal/models"
    "financial-service/internal/repository"
    "fmt"
)

type TransactionRepository struct {
    db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
    return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
    query := `
        INSERT INTO transactions 
        (from_user_id, to_user_id, amount, type, status, created_at)
        VALUES 
        (NULLIF(?, 0), NULLIF(?, 0), ?, ?, ?, ?)
    `
    
    result, err := r.db.ExecContext(ctx, query,
        tx.FromUserID,
        tx.ToUserID,
        tx.Amount,
        tx.Type,
        tx.Status,
        tx.CreatedAt,
    )
    if err != nil {
        return fmt.Errorf("failed to create transaction: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get last insert id: %w", err)
    }

    tx.ID = uint(id)
    return nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uint) (*models.Transaction, error) {
    tx := &models.Transaction{}
    query := `
        SELECT id, from_user_id, to_user_id, amount, type, status, created_at
        FROM transactions WHERE id = ?
    `
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &tx.ID,
        &tx.FromUserID,
        &tx.ToUserID,
        &tx.Amount,
        &tx.Type,
        &tx.Status,
        &tx.CreatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, repository.ErrNotFound
    }
    if err != nil {
        return nil, err
    }
    return tx, nil
}

func (r *TransactionRepository) UpdateStatus(ctx context.Context, id uint, status models.TransactionStatus) error {
    query := `UPDATE transactions SET status = ? WHERE id = ?`
    result, err := r.db.ExecContext(ctx, query, status, id)
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

func (r *TransactionRepository) GetUserTransactions(ctx context.Context, userID uint, limit, offset int) ([]models.Transaction, error) {
    query := `
        SELECT id, from_user_id, to_user_id, amount, type, status, created_at
        FROM transactions 
        WHERE from_user_id = ? OR to_user_id = ?
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `
    rows, err := r.db.QueryContext(ctx, query, userID, userID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []models.Transaction
    for rows.Next() {
        var tx models.Transaction
        err := rows.Scan(
            &tx.ID,
            &tx.FromUserID,
            &tx.ToUserID,
            &tx.Amount,
            &tx.Type,
            &tx.Status,
            &tx.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, tx)
    }
    return transactions, nil
} 