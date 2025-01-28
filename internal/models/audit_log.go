package models

import (
    "time"
    "encoding/json"
)

type AuditLog struct {
    ID         uint            `json:"id"`
    EntityType string         `json:"entity_type"`
    EntityID   uint           `json:"entity_id"`
    Action     string         `json:"action"`
    Details    json.RawMessage `json:"details"`
    Changes    string         `json:"changes"`
    CreatedAt  time.Time      `json:"created_at"`
} 