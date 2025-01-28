package models

import (
    "sync"
    "time"
    "errors"
)

type Balance struct {
    mu            sync.RWMutex `json:"-"`
    UserID        uint      `json:"user_id"`
    Amount        float64   `json:"amount"`
    LastUpdatedAt time.Time `json:"last_updated_at"`
}

func (b *Balance) GetAmount() float64 {
    b.mu.RLock()
    defer b.mu.RUnlock()

    return b.Amount
}

func (b *Balance) UpdateAmount(amount float64) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.Amount = amount
    b.LastUpdatedAt = time.Now()
}

func (b *Balance) AddAmount(amount float64) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.Amount += amount
    b.LastUpdatedAt = time.Now()
}

func (b *Balance) SubtractAmount(amount float64) error {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.Amount < amount {
        return errors.New("insufficient balance")
    }

    b.Amount -= amount
    b.LastUpdatedAt = time.Now()
    
    return nil
} 