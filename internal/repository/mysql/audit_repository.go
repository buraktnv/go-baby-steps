package mysql

import (
    "context"
    "database/sql"
    "financial-service/internal/models"
)

type AuditLogRepository struct {
    db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
    return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {
    query := `
        INSERT INTO audit_logs (
            entity_type, entity_id, action, changes, created_at
        ) VALUES (?, ?, ?, ?, ?)
    `
    _, err := r.db.ExecContext(ctx, query,
        log.EntityType,
        log.EntityID,
        log.Action,
        log.Changes,
        log.CreatedAt,
    )

    return err
}

func (r *AuditLogRepository) GetByEntityID(ctx context.Context, entityType string, entityID uint) ([]*models.AuditLog, error) {
    query := `
        SELECT entity_type, entity_id, action, changes, created_at
        FROM audit_logs 
        WHERE entity_type = ? AND entity_id = ?
        ORDER BY created_at DESC
    `

    rows, err := r.db.QueryContext(ctx, query, entityType, entityID)

    if err != nil {
        return nil, err
    }

    defer rows.Close()

    var logs []*models.AuditLog

    for rows.Next() {
        log := &models.AuditLog{}

        err := rows.Scan(
            &log.EntityType,
            &log.EntityID,
            &log.Action,
            &log.Changes,
            &log.CreatedAt,
        )

        if err != nil {
            return nil, err
        }

        logs = append(logs, log)
    }
    
    return logs, nil
} 