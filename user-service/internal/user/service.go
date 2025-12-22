package user

import (
	"NotifyProject/user-service/internal/events"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type EventProducer interface {
	Produce(topic, key string, message []byte) error
}

type IRepository interface {
	Create(ctx context.Context, u *User) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}

type Service struct {
	repo     IRepository
	producer EventProducer
}

func NewService(repo IRepository, producer EventProducer) *Service {
	return &Service{
		repo:     repo,
		producer: producer,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) (*User, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u, err = s.repo.Create(
		ctx,
		&User{
			Email:    email,
			Password: string(hashedPassword),
		},
	)
	if err != nil {
		return nil, err
	}

	event := events.UserEvent{
		Type: "UserCreated",
		Payload: events.UserCreatedPayload{
			ID:    u.ID,
			Email: u.Email,
		},
	}

	message, _ := json.Marshal(event)
	_ = s.producer.Produce(events.TopicUserEvents, u.ID.String(), message)

	return u, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*User, error) {
	existedUser, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existedUser == nil {
		return nil, ErrWrongCredentials
	}
	if existedUser.DeletedAt != nil {
		return nil, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return nil, ErrWrongCredentials
	}

	return existedUser, nil
}

func (s *Service) Delete(ctx context.Context, userID uuid.UUID) error {
	err := s.repo.Delete(ctx, userID)
	if err != nil {
		return err
	}

	event := events.UserEvent{
		Type: "UserDeleted",
		Payload: events.UserDeletedPayload{
			ID: userID,
		},
	}

	message, _ := json.Marshal(event)
	_ = s.producer.Produce(events.TopicUserEvents, userID.String(), message)

	return nil
}
