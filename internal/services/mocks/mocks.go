package mocks

import (
    "context"
    "financial-service/internal/models"
    "github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

type MockBalanceRepository struct {
    mock.Mock
}

func (m *MockBalanceRepository) GetBalance(ctx context.Context, userID uint) (*models.Balance, error) {
    args := m.Called(ctx, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Balance), args.Error(1)
}

func (m *MockBalanceRepository) UpdateBalance(ctx context.Context, balance *models.Balance) error {
    args := m.Called(ctx, balance)
    return args.Error(0)
}

func (m *MockBalanceRepository) CreateBalance(ctx context.Context, balance *models.Balance) error {
    args := m.Called(ctx, balance)
    return args.Error(0)
}

type MockTransactionRepository struct {
    mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
    args := m.Called(ctx, tx)
    return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id uint) (*models.Transaction, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) UpdateStatus(ctx context.Context, id uint, status models.TransactionStatus) error {
    args := m.Called(ctx, id, status)
    return args.Error(0)
}

func (m *MockTransactionRepository) GetUserTransactions(ctx context.Context, userID uint, limit, offset int) ([]models.Transaction, error) {
    args := m.Called(ctx, userID, limit, offset)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]models.Transaction), args.Error(1)
}

type MockAuditLogRepository struct {
    mock.Mock
}

func (m *MockAuditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {
    args := m.Called(ctx, log)
    return args.Error(0)
}

func (m *MockAuditLogRepository) GetByEntityID(ctx context.Context, entityType string, entityID uint) ([]*models.AuditLog, error) {
    args := m.Called(ctx, entityType, entityID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*models.AuditLog), args.Error(1)
} 