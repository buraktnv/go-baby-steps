package services

import (
    "context"
    "errors"
    "time"
    "financial-service/internal/models"
    "financial-service/internal/repository"
)

type UserService struct {
    userRepo    repository.UserRepository
    balanceRepo repository.BalanceRepository
    auditLogger *AuditLogger
}

func NewUserService(userRepo repository.UserRepository, balanceRepo repository.BalanceRepository) *UserService {
    return &UserService{
        userRepo:    userRepo,
        balanceRepo: balanceRepo,
    }
}

func (s *UserService) SetAuditLogger(logger *AuditLogger) {
    s.auditLogger = logger
}

// RegisterUser creates a new user with initial balance
func (s *UserService) RegisterUser(ctx context.Context, username, email, password string) (*models.User, error) {
    // Create user
    user := &models.User{
        Username:  username,
        Email:     email,
        Role:      models.RoleUser,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // Validate user data
    if err := user.Validate(); err != nil {
        return nil, err
    }

    // Hash password
    if err := user.SetPassword(password); err != nil {
        return nil, err
    }

    // Save user
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    // Create initial balance
    balance := &models.Balance{
        UserID:        user.ID,
        Amount:        0,
        LastUpdatedAt: time.Now(),
    }
    if err := s.balanceRepo.CreateBalance(ctx, balance); err != nil {
        return nil, err
    }

    return user, nil
}

// AuthenticateUser verifies user credentials and returns a user if valid
func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    if !user.CheckPassword(password) {
        return nil, errors.New("invalid credentials")
    }

    return user, nil
}

// HasPermission checks if user has required role
func (s *UserService) HasPermission(ctx context.Context, userID uint, requiredRole models.Role) bool {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return false
    }
    return user.Role == requiredRole || user.Role == models.RoleAdmin
} 