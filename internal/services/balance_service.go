package services

import (
    "context"
    "sync"
    "time"
    "financial-service/internal/models"
    "financial-service/internal/repository"
)

type BalanceService struct {
    balanceRepo repository.BalanceRepository
    txRepo      repository.TransactionRepository
    cache       *BalanceCache
}

type BalanceCache struct {
    balances map[uint]*models.Balance
    mu       sync.RWMutex
}

func NewBalanceService(balanceRepo repository.BalanceRepository, txRepo repository.TransactionRepository) *BalanceService {
    return &BalanceService{
        balanceRepo: balanceRepo,
        txRepo:      txRepo,
        cache: &BalanceCache{
            balances: make(map[uint]*models.Balance),
        },
    }
}

func (s *BalanceService) GetBalance(ctx context.Context, userID uint) (*models.Balance, error) {
    if balance := s.cache.get(userID); balance != nil {
        return balance, nil
    }

    balance, err := s.balanceRepo.GetBalance(ctx, userID)
    if err != nil {
        return nil, err
    }

    s.cache.set(userID, balance)
    return balance, nil
}

func (s *BalanceService) RecalculateBalance(ctx context.Context, userID uint) error {
    transactions, err := s.txRepo.GetUserTransactions(ctx, userID, 0, 0)
    if err != nil {
        return err
    }

    var totalBalance float64
    for _, tx := range transactions {
        if tx.Status != models.TransactionStatusCompleted {
            continue
        }

        if tx.ToUserID == userID {
            totalBalance += tx.Amount
        }
        if tx.FromUserID == userID {
            totalBalance -= tx.Amount
        }
    }

    balance := &models.Balance{
        UserID:        userID,
        Amount:        totalBalance,
        LastUpdatedAt: time.Now(),
    }

    if err := s.balanceRepo.UpdateBalance(ctx, balance); err != nil {
        return err
    }

    s.cache.set(userID, balance)
    return nil
}

func (c *BalanceCache) get(userID uint) *models.Balance {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.balances[userID]
}

func (c *BalanceCache) set(userID uint, balance *models.Balance) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.balances[userID] = balance
} 