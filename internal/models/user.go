package models

import (
    "errors"
    "time"
    "golang.org/x/crypto/bcrypt"
)

type Role string

const (
    RoleUser  Role = "user"
    RoleAdmin Role = "admin"
)

type User struct {
    ID           uint      `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    Role         Role      `json:"role"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// Validate checks if user data is valid
func (u *User) Validate() error {
    if u.Username == "" {
        return errors.New("username is required")
    }

    if len(u.Username) < 3 {
        return errors.New("username must be at least 3 characters")
    }

    if u.Email == "" {
        return errors.New("email is required")
    }
    // Add more validation as needed
    return nil
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
    if len(password) < 6 {
        return errors.New("password must be at least 6 characters")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

    if err != nil {
        return err
    }

    u.PasswordHash = string(hash)

    return nil
}

// CheckPassword verifies the user's password
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
    
    return err == nil
} 