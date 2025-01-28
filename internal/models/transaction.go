package models

import (
    "errors"
    "sync"
    "time"
)

type TransactionType string
type TransactionStatus string

const (
    TransactionTypeCredit   TransactionType = "credit"
    TransactionTypeDebit    TransactionType = "debit"
    TransactionTypeTransfer TransactionType = "transfer"

    TransactionStatusPending   TransactionStatus = "pending"
    TransactionStatusCompleted TransactionStatus = "completed"
    TransactionStatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
    mu          sync.RWMutex     `json:"-"`
    ID          uint             `json:"id"`
    FromUserID  uint             `json:"from_user_id"`
    ToUserID    uint             `json:"to_user_id"`
    Amount      float64          `json:"amount"`
    Type        TransactionType  `json:"type"`
    Status      TransactionStatus `json:"status"`
    CreatedAt   time.Time        `json:"created_at"`
}

func (t *Transaction) SetStatus(status TransactionStatus) {
    t.mu.Lock()
    defer t.mu.Unlock()

    t.Status = status
}

func (t *Transaction) GetStatus() TransactionStatus {
    t.mu.RLock()
    defer t.mu.RUnlock()

    return t.Status
}

func (t *Transaction) Validate() error {
    switch t.Type {
        case TransactionTypeCredit:
            if t.ToUserID == 0 {
                return errors.New("to_user_id is required for credit transactions")
            }
        case TransactionTypeDebit:
            if t.FromUserID == 0 {
                return errors.New("from_user_id is required for debit transactions")
            }
        case TransactionTypeTransfer:
            if t.FromUserID == 0 || t.ToUserID == 0 {
                return errors.New("both from_user_id and to_user_id are required for transfers")
            }
        default:
            return errors.New("invalid transaction type")
    }

    if t.Amount <= 0 {
        return errors.New("amount must be positive")
    }

    return nil
} 