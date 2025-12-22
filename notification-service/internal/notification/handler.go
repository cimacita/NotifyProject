package notification

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"NotifyProject/notification-service/pkg/auth"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ServiceInterface interface {
	Send(ctx context.Context, sender, reciever uuid.UUID, message string) (*Notification, error)
	Inbox(ctx context.Context, receiver uuid.UUID) ([]Notification, error)
	UserIdByNotif(ctx context.Context, notifID uuid.UUID) (uuid.UUID, error)
	Read(ctx context.Context, notifID uuid.UUID) error
	Remove(ctx context.Context, notifID uuid.UUID) error
}

type Handler struct {
	service    ServiceInterface
	jwtManager *auth.JWTManager
}

func NewHandler(router *http.ServeMux, s ServiceInterface, jwtManager *auth.JWTManager) {
	handler := &Handler{
		service:    s,
		jwtManager: jwtManager,
	}
	router.Handle("POST /notifications", auth.Middleware(jwtManager, http.HandlerFunc(handler.SendNotification)))
	router.Handle("GET /notifications", auth.Middleware(jwtManager, http.HandlerFunc(handler.GetAllNotifications)))
	router.Handle("PATCH /notifications/{id}", auth.Middleware(jwtManager, http.HandlerFunc(handler.ReadNotification)))
	router.Handle("DELETE /notifications/{id}", auth.Middleware(jwtManager, http.HandlerFunc(handler.DeleteNotification)))
}

func (handler *Handler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var payload SendNotificationRequest

	err := decodeJSON[SendNotificationRequest](r, &payload)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateStruct[SendNotificationRequest](payload)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	sender, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, "", http.StatusUnauthorized)
		return
	}

	n, err := handler.service.Send(r.Context(), sender, payload.Receiver, payload.Message)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := SendNotificationResponse{
		ID:        n.ID,
		Sender:    n.Sender,
		Receiver:  n.Receiver,
		Message:   n.Message,
		CreatedAt: time.Now(),
	}

	writeJSON(w, data, http.StatusCreated)
}

func (handler *Handler) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	receiver, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, "", http.StatusUnauthorized)
		return
	}

	ns, err := handler.service.Inbox(r.Context(), receiver)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var inbox InboxResponse

	for _, notif := range ns {
		n := NotificationResponse{
			ID:        notif.ID,
			Sender:    notif.Sender,
			Receiver:  notif.Receiver,
			Message:   notif.Message,
			CreatedAt: notif.CreatedAt,
			ReadAt:    notif.ReadAt,
		}

		inbox.Notifications = append(inbox.Notifications, n)
	}

	writeJSON(w, inbox, http.StatusOK)
}

func (handler *Handler) ReadNotification(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, "", http.StatusUnauthorized)
		return
	}

	notifID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	userByNotif, err := handler.service.UserIdByNotif(r.Context(), notifID)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userID != userByNotif {
		writeJSON(w, "", http.StatusUnauthorized)
		return
	}

	err = handler.service.Read(r.Context(), notifID)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, "marked as read", http.StatusOK)
}

func (handler *Handler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, "", http.StatusUnauthorized)
		return
	}

	notifID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	userByNotif, err := handler.service.UserIdByNotif(r.Context(), notifID)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userID != userByNotif {
		writeJSON(w, "", http.StatusUnauthorized)
		return
	}

	err = handler.service.Remove(r.Context(), notifID)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, "notification deleted", http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func decodeJSON[T any](r *http.Request, v *T) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func validateStruct[T any](v T) error {
	return validator.New().Struct(v)
}
