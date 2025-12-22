package userShadow

import (
	"github.com/google/uuid"
)

type UserShadow struct {
	ID        uuid.UUID
	IsDeleted bool
}
