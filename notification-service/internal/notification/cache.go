package notification

import (
	"NotifyProject/notification-service/pkg/cache"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	inboxKeyPrefix = "inbox:user:"
	inboxTTL       = 5 * time.Minute
)

type NotifCache struct {
	cache cache.Cache[[]Notification]
}

func NewNotifCache(cache cache.Cache[[]Notification]) *NotifCache {
	return &NotifCache{cache: cache}
}

func (n *NotifCache) GetInbox(ctx context.Context, receiver uuid.UUID) ([]Notification, error) {
	key := fmt.Sprintf("%s%s", inboxKeyPrefix, receiver.String())
	return n.cache.Get(ctx, key)
}

func (n *NotifCache) SetInbox(ctx context.Context, receiver uuid.UUID, notifs []Notification) error {
	key := fmt.Sprintf("%s%s", inboxKeyPrefix, receiver.String())
	return n.cache.Set(ctx, key, notifs, inboxTTL)
}

func (n *NotifCache) DeleteInbox(ctx context.Context, receiver uuid.UUID) error {
	key := fmt.Sprintf("%s%s", inboxKeyPrefix, receiver.String())
	return n.cache.Delete(ctx, key)
}
