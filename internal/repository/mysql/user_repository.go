package mysql

import (
    "context"
    "database/sql"
    "financial-service/internal/models"
    "financial-service/internal/repository"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
    result, err := r.db.ExecContext(ctx, query,
        user.Username,
        user.Email,
        user.PasswordHash,
        user.Role,
        user.CreatedAt,
        user.UpdatedAt,
    )
    if err != nil {
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return err
    }

    user.ID = uint(id)
    return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
    user := &models.User{}
    query := `
        SELECT id, username, email, password_hash, role, created_at, updated_at
        FROM users WHERE id = ?
    `
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, repository.ErrNotFound
    }
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    user := &models.User{}
    query := `
        SELECT id, username, email, password_hash, role, created_at, updated_at
        FROM users WHERE email = ?
    `
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, repository.ErrNotFound
    }
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
    query := `
        UPDATE users 
        SET username = ?, email = ?, password_hash = ?, role = ?, updated_at = ?
        WHERE id = ?
    `
    result, err := r.db.ExecContext(ctx, query,
        user.Username,
        user.Email,
        user.PasswordHash,
        user.Role,
        user.UpdatedAt,
        user.ID,
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