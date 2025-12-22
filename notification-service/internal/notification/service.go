package notification

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

var SystemUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type INotifRepository interface {
	Create(ctx context.Context, n *Notification) (*Notification, error)
	GetByReceiver(ctx context.Context, receiver uuid.UUID) ([]Notification, error)
	GetUserByNotifID(ctx context.Context, notifID uuid.UUID) (uuid.UUID, error)
	MarkRead(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteAllForUser(ctx context.Context, receiver uuid.UUID) error
}

type INotifCache interface {
	GetInbox(ctx context.Context, receiver uuid.UUID) ([]Notification, error)
	SetInbox(ctx context.Context, receiver uuid.UUID, notifs []Notification) error
	DeleteInbox(ctx context.Context, receiver uuid.UUID) error
}

type IUserShadowRepository interface {
	IsDeleted(ctx context.Context, id uuid.UUID) (bool, error)
}

type Service struct {
	notifRepo      INotifRepository
	notifCache     INotifCache
	userShadowRepo IUserShadowRepository
}

func NewService(notifRepo INotifRepository, notifCache INotifCache, userShadowRepo IUserShadowRepository) *Service {
	return &Service{
		notifRepo:      notifRepo,
		notifCache:     notifCache,
		userShadowRepo: userShadowRepo,
	}
}

func (s *Service) Send(ctx context.Context, sender, receiver uuid.UUID, message string) (*Notification, error) {
	if sender != SystemUserID {
		if err := s.checkUserActive(ctx, sender); err != nil {
			return nil, fmt.Errorf("sender check failed: %w", err)
		}
	}
	if err := s.checkUserActive(ctx, receiver); err != nil {
		return nil, fmt.Errorf("receiver check failed: %w", err)
	}

	n, err := s.notifRepo.Create(
		ctx,
		&Notification{
			Sender:   sender,
			Receiver: receiver,
			Message:  message,
		},
	)
	if err != nil {
		return nil, err
	}

	_ = s.notifCache.DeleteInbox(ctx, receiver)

	return n, nil
}

func (s *Service) SendWelcomeNotif(ctx context.Context, receiver uuid.UUID, email string) error {
	message := fmt.Sprintf("Welcome to our service! Your registered email is: %s", email)

	_, err := s.Send(ctx, SystemUserID, receiver, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Inbox(ctx context.Context, receiver uuid.UUID) ([]Notification, error) {
	if err := s.checkUserActive(ctx, receiver); err != nil {
		return nil, err
	}

	notifications, err := s.notifCache.GetInbox(ctx, receiver)
	if err == nil {
		return notifications, nil
	} else {
		log.Printf("WARNING: failed to get redis cache: %v", err)
	}

	notifications, err = s.notifRepo.GetByReceiver(ctx, receiver)
	if err != nil {
		return nil, err
	}

	err = s.notifCache.SetInbox(ctx, receiver, notifications)
	if err != nil {
		log.Printf("WARNING: failed to set redis cache: %v", err)
	}

	return notifications, nil
}

func (s *Service) UserIdByNotif(ctx context.Context, notifID uuid.UUID) (uuid.UUID, error) {
	userID, err := s.notifRepo.GetUserByNotifID(ctx, notifID)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, err
}

func (s *Service) Read(ctx context.Context, notifID uuid.UUID) error {
	userID, err := s.UserIdByNotif(ctx, notifID)
	if err != nil {
		return err
	}
	if err := s.checkUserActive(ctx, userID); err != nil {
		return err
	}

	err = s.notifRepo.MarkRead(ctx, notifID)
	if err != nil {
		return err
	}

	_ = s.notifCache.DeleteInbox(ctx, userID)

	return nil
}

func (s *Service) Remove(ctx context.Context, notifID uuid.UUID) error {
	userID, err := s.UserIdByNotif(ctx, notifID)
	if err != nil {
		return err
	}
	if err := s.checkUserActive(ctx, userID); err != nil {
		return err
	}

	err = s.notifRepo.Delete(ctx, notifID)
	if err != nil {
		return err
	}

	_ = s.notifCache.DeleteInbox(ctx, userID)

	return nil
}

func (s *Service) RemoveNotifsForUser(ctx context.Context, receiver uuid.UUID) error {
	err := s.notifRepo.DeleteAllForUser(ctx, receiver)
	if err != nil {
		return err
	}

	_ = s.notifCache.DeleteInbox(ctx, receiver)

	return nil
}

func (s *Service) checkUserActive(ctx context.Context, userID uuid.UUID) error {
	IsDeleted, err := s.userShadowRepo.IsDeleted(ctx, userID)
	if err != nil {
		return err
	}

	if IsDeleted {
		return fmt.Errorf("user %s is deleted", userID)
	}

	return nil
}
