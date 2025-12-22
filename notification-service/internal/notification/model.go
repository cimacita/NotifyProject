package notification

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID
	Sender    uuid.UUID
	Receiver  uuid.UUID
	Message   string
	CreatedAt time.Time
	ReadAt    *time.Time
}
