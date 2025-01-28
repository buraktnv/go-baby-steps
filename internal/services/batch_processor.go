package services

import (
    "context"
    "fmt"
    "sync"
    "financial-service/internal/models"
    "financial-service/internal/repository"
    "github.com/rs/zerolog/log"
)

type BatchProcessor struct {
    txRepo      repository.TransactionRepository
    balanceRepo repository.BalanceRepository
    workerPool  *WorkerPool
    batchSize   int
    mu          sync.Mutex
    processing  bool
}

func NewBatchProcessor(
    txRepo repository.TransactionRepository,
    balanceRepo repository.BalanceRepository,
    workerPool *WorkerPool,
    batchSize int,
) *BatchProcessor {
    return &BatchProcessor{
        txRepo:      txRepo,
        balanceRepo: balanceRepo,
        workerPool:  workerPool,
        batchSize:   batchSize,
    }
}

func (bp *BatchProcessor) ProcessPendingTransactions(ctx context.Context) error {
    bp.mu.Lock()
    if bp.processing {
        bp.mu.Unlock()
        return fmt.Errorf("batch processing already in progress")
    }
    bp.processing = true
    bp.mu.Unlock()

    defer func() {
        bp.mu.Lock()
        bp.processing = false
        bp.mu.Unlock()
    }()

    // Get pending transactions
    transactions, err := bp.getPendingTransactions(ctx)
    if err != nil {
        return fmt.Errorf("failed to get pending transactions: %w", err)
    }

    if len(transactions) == 0 {
        return nil
    }

    // Process transactions in parallel
    errChan := make(chan error, len(transactions))
    for _, tx := range transactions {
        resultChan := make(chan error, 1)
        err := bp.workerPool.Submit(&Task{
            Transaction: tx,
            ResultChan:  resultChan,
        })
        if err != nil {
            log.Error().Err(err).Uint("tx_id", tx.ID).Msg("Failed to submit transaction")
            continue
        }

        go func(tx *models.Transaction) {
            if err := <-resultChan; err != nil {
                errChan <- fmt.Errorf("failed to process transaction %d: %w", tx.ID, err)
            } else {
                errChan <- nil
            }
        }(tx)
    }

    // Collect errors
    var errors []error
    for i := 0; i < len(transactions); i++ {
        if err := <-errChan; err != nil {
            errors = append(errors, err)
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("batch processing completed with %d errors: %v", len(errors), errors)
    }

    return nil
}

func (bp *BatchProcessor) getPendingTransactions(ctx context.Context) ([]*models.Transaction, error) {
    // Implementation depends on your repository
    // This is just a placeholder
    return nil, nil
} 