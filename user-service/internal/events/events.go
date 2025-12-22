package events

import "github.com/google/uuid"

const (
	TopicUserEvents = "user-events"
)

type UserEvent struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type UserCreatedPayload struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type UserDeletedPayload struct {
	ID uuid.UUID `json:"id"`
}
