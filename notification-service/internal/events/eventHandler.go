package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

type IRepository interface {
	InsertIfNotExists(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type IService interface {
	SendWelcomeNotif(ctx context.Context, receiver uuid.UUID, email string) error
	RemoveNotifsForUser(ctx context.Context, receiver uuid.UUID) error
}

type UserEventHandler struct {
	repo    IRepository
	service IService
}

func NewUserEventHandler(repo IRepository, service IService) *UserEventHandler {
	return &UserEventHandler{
		repo:    repo,
		service: service,
	}
}

func (handler *UserEventHandler) Handle(msg *kafka.Message) error {
	var event UserEvent
	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		log.Printf("ERROR: Failed to unmarshal message body (Type/RawPayload): %v. Raw: %s", err, string(msg.Value))
		return nil
	}

	ctx := context.Background()
	switch event.Type {

	case UserCreated:
		var payload UserCreatedPayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			log.Printf("ERROR: Failed to unmarshal UserCreated payload: %v", err)
			return nil
		}

		return handler.handleUserCreated(ctx, payload)

	case UserDeleted:
		var payload UserDeletedPayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			log.Printf("ERROR: Failed to unmarshal UserDeleted payload: %v", err)
			return nil
		}

		return handler.handleUserDeleted(ctx, payload)

	default:
		log.Printf("WARNING: Unknown event type: %s", event.Type)
		return nil
	}
}

func (handler *UserEventHandler) handleUserCreated(ctx context.Context, payload UserCreatedPayload) error {
	err := handler.repo.InsertIfNotExists(ctx, payload.ID)
	if err != nil {
		return fmt.Errorf("failed to sync user creation: %w", err)
	}

	err = handler.service.SendWelcomeNotif(ctx, payload.ID, payload.Email)
	if err != nil {
		return fmt.Errorf("failed to send welcome notification: %w", err)
	}

	return nil
}

func (handler *UserEventHandler) handleUserDeleted(ctx context.Context, payload UserDeletedPayload) error {
	err := handler.repo.Delete(ctx, payload.ID)
	if err != nil {
		return fmt.Errorf("failed to sync user deletion: %w", err)
	}

	err = handler.service.RemoveNotifsForUser(ctx, payload.ID)
	if err != nil {
		return fmt.Errorf("failed to remove notifs for user: %w", err)
	}

	return nil
}
