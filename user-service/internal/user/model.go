package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
	DeletedAt *time.Time
}
