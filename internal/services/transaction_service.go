package services

import (
    "context"
    "errors"
    "fmt"
    "financial-service/internal/models"
    "financial-service/internal/repository"
    "time"
    "github.com/rs/zerolog/log"
)

type TransactionService struct {
    txRepo      repository.TransactionRepository
    balanceRepo repository.BalanceRepository
    userRepo    repository.UserRepository
    workerPool  *WorkerPool
    auditLogger *AuditLogger
}

func NewTransactionService(
    txRepo repository.TransactionRepository,
    balanceRepo repository.BalanceRepository,
    userRepo repository.UserRepository,
    numWorkers int,
) *TransactionService {
    service := &TransactionService{
        txRepo:      txRepo,
        balanceRepo: balanceRepo,
        userRepo:    userRepo,
    }
    
    service.workerPool = NewWorkerPool(numWorkers, context.Background(), txRepo, balanceRepo, context.Background())
    service.workerPool.Start()
    
    return service
}

func (s *TransactionService) SetAuditLogger(logger *AuditLogger) {
    s.auditLogger = logger
}

func (s *TransactionService) Credit(ctx context.Context, userID uint, amount float64) (*models.Transaction, error) {
    // Validate user exists
    _, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        if err == repository.ErrNotFound {
            return nil, fmt.Errorf("user not found: %d", userID)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    tx := &models.Transaction{
        FromUserID: 0,
        ToUserID:   userID,
        Amount:     amount,
        Type:       models.TransactionTypeCredit,
        Status:     models.TransactionStatusPending,
        CreatedAt:  time.Now(),
    }

    if err := s.txRepo.Create(ctx, tx); err != nil {
        return nil, fmt.Errorf("failed to create transaction: %w", err)
    }

    // Process transaction
    resultChan := make(chan error, 1)
    err = s.workerPool.Submit(&Task{
        Transaction: tx,
        ResultChan:  resultChan,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to submit transaction: %w", err)
    }

    // Wait for processing
    if err := <-resultChan; err != nil {
        return nil, fmt.Errorf("failed to process transaction: %w", err)
    }

    // Log the audit
    if s.auditLogger != nil {
        changes := map[string]interface{}{
            "amount":    amount,
            "user_id":   userID,
            "type":      "credit",
            "status":    "completed",
        }
        if err := s.auditLogger.LogAction(ctx, "transaction", tx.ID, "credit", changes); err != nil {
            log.Error().Err(err).Msg("Failed to log audit")
        }
    }

    return tx, nil
}

func (s *TransactionService) Debit(ctx context.Context, userID uint, amount float64) (*models.Transaction, error) {
    // Validate user exists
    _, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        if err == repository.ErrNotFound {
            return nil, fmt.Errorf("user not found: %d", userID)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    // Validate balance
    balance, err := s.balanceRepo.GetBalance(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get balance: %w", err)
    }
    if balance.Amount < amount {
        return nil, errors.New("insufficient funds")
    }

    tx := &models.Transaction{
        FromUserID: userID,
        Amount:    amount,
        Type:      models.TransactionTypeDebit,
        Status:    models.TransactionStatusPending,
        CreatedAt: time.Now(),
    }

    if err := s.txRepo.Create(ctx, tx); err != nil {
        return nil, fmt.Errorf("failed to create transaction: %w", err)
    }

    // Process transaction
    resultChan := make(chan error, 1)
    err = s.workerPool.Submit(&Task{
        Transaction: tx,
        ResultChan:  resultChan,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to submit transaction: %w", err)
    }

    // Wait for processing
    if err := <-resultChan; err != nil {
        return nil, fmt.Errorf("failed to process transaction: %w", err)
    }

    // Log the audit
    if s.auditLogger != nil {
        changes := map[string]interface{}{
            "amount":    amount,
            "user_id":   userID,
            "type":      "debit",
            "status":    "completed",
        }
        if err := s.auditLogger.LogAction(ctx, "transaction", tx.ID, "debit", changes); err != nil {
            log.Error().Err(err).Msg("Failed to log audit")
        }
    }

    return tx, nil
}

func (s *TransactionService) Transfer(ctx context.Context, fromUserID, toUserID uint, amount float64) (*models.Transaction, error) {
    // Validate amount
    if amount <= 0 {
        return nil, errors.New("amount must be positive")
    }

    // Validate users exist
    _, err := s.userRepo.GetByID(ctx, fromUserID)
    if err != nil {
        if err == repository.ErrNotFound {
            return nil, fmt.Errorf("from user not found: %d", fromUserID)
        }
        return nil, fmt.Errorf("failed to get from user: %w", err)
    }

    _, err = s.userRepo.GetByID(ctx, toUserID)
    if err != nil {
        if err == repository.ErrNotFound {
            return nil, fmt.Errorf("to user not found: %d", toUserID)
        }
        return nil, fmt.Errorf("failed to get to user: %w", err)
    }

    // Validate balance
    balance, err := s.balanceRepo.GetBalance(ctx, fromUserID)
    if err != nil {
        if err != repository.ErrNotFound {
            return nil, fmt.Errorf("failed to get balance: %w", err)
        }
        // If balance not found, treat as zero
        balance = &models.Balance{
            UserID: fromUserID,
            Amount: 0,
        }
    }

    if balance.Amount < amount {
        return nil, errors.New("insufficient funds")
    }

    tx := &models.Transaction{
        FromUserID: fromUserID,
        ToUserID:   toUserID,
        Amount:     amount,
        Type:       models.TransactionTypeTransfer,
        Status:     models.TransactionStatusPending,
        CreatedAt:  time.Now(),
    }

    if err := s.txRepo.Create(ctx, tx); err != nil {
        return nil, fmt.Errorf("failed to create transaction: %w", err)
    }

    // Process transaction
    resultChan := make(chan error, 1)
    err = s.workerPool.Submit(&Task{
        Transaction: tx,
        ResultChan:  resultChan,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to submit transaction: %w", err)
    }

    // Wait for processing
    if err := <-resultChan; err != nil {
        return nil, fmt.Errorf("failed to process transaction: %w", err)
    }

    // Log the audit
    if s.auditLogger != nil {
        changes := map[string]interface{}{
            "amount":      amount,
            "from_user":   fromUserID,
            "to_user":     toUserID,
            "type":        "transfer",
            "status":      "completed",
        }
        if err := s.auditLogger.LogAction(ctx, "transaction", tx.ID, "transfer", changes); err != nil {
            log.Error().Err(err).Msg("Failed to log audit")
        }
    }

    return tx, nil
}

func (s *TransactionService) Cleanup() {
    if s.workerPool != nil {
        s.workerPool.Stop()
    }
} 