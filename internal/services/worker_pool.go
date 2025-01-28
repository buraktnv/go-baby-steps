package services

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
    "financial-service/internal/models"
    "financial-service/internal/repository"
)

type WorkerPool struct {
    numWorkers  int
    taskQueue   chan *Task
    wg          sync.WaitGroup
    ctx         context.Context
    cancel      context.CancelFunc
    txRepo      repository.TransactionRepository
    balanceRepo repository.BalanceRepository
    stats       *WorkerStats
}

type Task struct {
    Transaction *models.Transaction
    ResultChan  chan error
}

type WorkerStats struct {
    ProcessedCount   int64
    SuccessCount    int64
    ErrorCount      int64
}

func NewWorkerPool(numWorkers int, ctx context.Context, txRepo repository.TransactionRepository, balanceRepo repository.BalanceRepository, parentCtx context.Context) *WorkerPool {
    if ctx == nil {
        ctx = context.Background()
    }

    ctx, cancel := context.WithCancel(ctx)

    return &WorkerPool{
        numWorkers:  numWorkers,
        taskQueue:   make(chan *Task, numWorkers*2),
        ctx:         ctx,
        cancel:      cancel,
        txRepo:      txRepo,
        balanceRepo: balanceRepo,
        stats:       &WorkerStats{},
    }
}

func (wp *WorkerPool) Start() {
    wp.wg.Add(wp.numWorkers)
    
    for i := 0; i < wp.numWorkers; i++ {
        go wp.worker()
    }
}

func (wp *WorkerPool) Stop() {
    wp.cancel()
    close(wp.taskQueue)
    wp.wg.Wait()
}

func (wp *WorkerPool) Submit(task *Task) error {
    select {
        case wp.taskQueue <- task:
            return nil
        case <-wp.ctx.Done():
            return errors.New("worker pool is stopped")
        default:
            return errors.New("task queue is full")
    }
}

func (wp *WorkerPool) worker() {
    defer wp.wg.Done()

    for {
        select {
            case <-wp.ctx.Done():
                return
            case task, ok := <-wp.taskQueue:
                if !ok {
                    return
                }

                err := wp.processTransaction(task.Transaction)
                atomic.AddInt64(&wp.stats.ProcessedCount, 1)

                if err != nil {
                    atomic.AddInt64(&wp.stats.ErrorCount, 1)
                } else {
                    atomic.AddInt64(&wp.stats.SuccessCount, 1)
                }

                task.ResultChan <- err
        }
    }
}

func (wp *WorkerPool) processTransaction(tx *models.Transaction) error {
    ctx, cancel := context.WithTimeout(wp.ctx, 5*time.Second)

    defer cancel()

    switch tx.Type {
        case models.TransactionTypeTransfer:
            // Debit from source account
            fromBalance, err := wp.balanceRepo.GetBalance(ctx, tx.FromUserID)

            if err != nil {
                return fmt.Errorf("failed to get source balance: %w", err)
            }
            
            if fromBalance.Amount < tx.Amount {
                return errors.New("insufficient funds")
            }
            
            fromBalance.Amount -= tx.Amount

            if err := wp.balanceRepo.UpdateBalance(ctx, fromBalance); err != nil {
                return fmt.Errorf("failed to update source balance: %w", err)
            }

            // Credit to destination account
            toBalance, err := wp.balanceRepo.GetBalance(ctx, tx.ToUserID)
            if err != nil {
                if err == repository.ErrNotFound {
                    // Create new balance if it doesn't exist
                    toBalance = &models.Balance{
                        UserID: tx.ToUserID,
                        Amount: 0,
                    }
                    if err := wp.balanceRepo.CreateBalance(ctx, toBalance); err != nil {
                        return fmt.Errorf("failed to create destination balance: %w", err)
                    }
                } else {
                    return fmt.Errorf("failed to get destination balance: %w", err)
                }
            }

            toBalance.Amount += tx.Amount
            if err := wp.balanceRepo.UpdateBalance(ctx, toBalance); err != nil {
                return fmt.Errorf("failed to update destination balance: %w", err)
            }

        case models.TransactionTypeCredit:
            balance, err := wp.balanceRepo.GetBalance(ctx, tx.ToUserID)

            if err != nil {
                if err == repository.ErrNotFound {
                    balance = &models.Balance{
                        UserID: tx.ToUserID,
                        Amount: 0,
                    }
                    if err := wp.balanceRepo.CreateBalance(ctx, balance); err != nil {
                        return fmt.Errorf("failed to create balance: %w", err)
                    }
                } else {
                    return fmt.Errorf("failed to get balance: %w", err)
                }
            }
            
            balance.Amount += tx.Amount

            if err := wp.balanceRepo.UpdateBalance(ctx, balance); err != nil {
                return fmt.Errorf("failed to update balance: %w", err)
            }

        case models.TransactionTypeDebit:
            balance, err := wp.balanceRepo.GetBalance(ctx, tx.FromUserID)

            if err != nil {
                return fmt.Errorf("failed to get balance: %w", err)
            }
            
            if balance.Amount < tx.Amount {
                return errors.New("insufficient funds")
            }
            
            balance.Amount -= tx.Amount
            
            if err := wp.balanceRepo.UpdateBalance(ctx, balance); err != nil {
                return fmt.Errorf("failed to update balance: %w", err)
            }
    }

    // Update transaction status
    if err := wp.txRepo.UpdateStatus(ctx, tx.ID, models.TransactionStatusCompleted); err != nil {
        return fmt.Errorf("failed to update transaction status: %w", err)
    }

    return nil
}

func (wp *WorkerPool) GetStats() WorkerStats {
    return WorkerStats{
        ProcessedCount: atomic.LoadInt64(&wp.stats.ProcessedCount),
        SuccessCount:  atomic.LoadInt64(&wp.stats.SuccessCount),
        ErrorCount:    atomic.LoadInt64(&wp.stats.ErrorCount),
    }
} 