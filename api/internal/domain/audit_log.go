package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// AuditLog represents an immutable, append only security log record
// Records every mutation in the system with actor, action, and JSON changes payload
type AuditLog struct {
	ID             int64           `json:"id" db:"id"`
	OccurredAt     time.Time       `json:"occurred_at" db:"occurred_at"`
	ActorID        string          `json:"actor_id,omitempty" db:"actor_id"`
	Action         string          `json:"action" db:"action"`
	EntityName     string          `json:"entity_name" db:"entity_name"`
	EntityID       string          `json:"entity_id,omitempty" db:"entity_id"`
	ChangesPayload json.RawMessage `json:"changes_payload,omitempty" db:"changes_payload"`
	ClientIP       string          `json:"client_ip,omitempty" db:"client_ip"`
}

// NewAuditLog creates and validates a new AuditLog record (Factory Method pattern)
func NewAuditLog(actorID, action, entityName, entityID string, changes interface{}, clientIP string) (*AuditLog, error) {
	action = strings.TrimSpace(action)
	if action == "" {
		return nil, errors.New("audit action cannot be empty")
	}

	entityName = strings.TrimSpace(entityName)
	if entityName == "" {
		return nil, errors.New("audit entity_name cannot be empty")
	}

	var rawChanges json.RawMessage
	if changes != nil {
		bytes, err := json.Marshal(changes)
		if err != nil {
			return nil, err
		}
		rawChanges = bytes
	}

	return &AuditLog{
		OccurredAt:     time.Now().UTC(),
		ActorID:        strings.TrimSpace(actorID),
		Action:         action,
		EntityName:     entityName,
		EntityID:       strings.TrimSpace(entityID),
		ChangesPayload: rawChanges,
		ClientIP:       strings.TrimSpace(clientIP),
	}, nil
}
