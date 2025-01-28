package services

import (
    "context"
    "encoding/json"
    "financial-service/internal/models"
    "financial-service/internal/repository"
    "time"
)

type AuditLogger struct {
    repo repository.AuditLogRepository
}

func NewAuditLogger(repo repository.AuditLogRepository) *AuditLogger {
    return &AuditLogger{
        repo: repo,
    }
}

func (l *AuditLogger) LogAction(ctx context.Context, entityType string, entityID uint, action string, changes interface{}) error {
    changesJSON, err := json.Marshal(changes)
    if err != nil {
        return err
    }

    log := &models.AuditLog{
        EntityType: entityType,
        EntityID:   entityID,
        Action:     action,
        Changes:    string(changesJSON),
        CreatedAt:  time.Now(),
    }

    return l.repo.Create(ctx, log)
} 