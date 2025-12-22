package events

import (
	"encoding/json"

	"github.com/google/uuid"
)

const (
	UserCreated = "UserCreated"
	UserDeleted = "UserDeleted"
)

type UserEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type UserCreatedPayload struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type UserDeletedPayload struct {
	ID uuid.UUID `json:"id"`
}
