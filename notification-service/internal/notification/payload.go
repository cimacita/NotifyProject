package notification

import (
	"time"

	"github.com/google/uuid"
)

type SendNotificationRequest struct {
	Receiver uuid.UUID `json:"receiver" validate:"required,uuid4"`
	Message  string    `json:"message" validate:"required,min=1"`
}

type SendNotificationResponse struct {
	ID        uuid.UUID `json:"id"`
	Sender    uuid.UUID `json:"sender"`
	Receiver  uuid.UUID `json:"receiver"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationResponse struct {
	ID        uuid.UUID  `json:"id"`
	Sender    uuid.UUID  `json:"sender"`
	Receiver  uuid.UUID  `json:"receiver"`
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at"`
}

type InboxResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
}

type ReadNotificationResponse struct {
	ID     uuid.UUID  `json:"id"`
	ReadAt *time.Time `json:"read_at"`
}
