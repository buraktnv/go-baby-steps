package services

type AuditLoggable interface {
    SetAuditLogger(logger *AuditLogger)
}