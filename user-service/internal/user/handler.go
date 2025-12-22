package user

import (
	"NotifyProject/user-service/pkg/auth"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ServiceInterface interface {
	Register(ctx context.Context, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}

type Handler struct {
	service    ServiceInterface
	jwtManager *auth.JWTManager
}

func NewHandler(router *http.ServeMux, service ServiceInterface, jwtManager *auth.JWTManager) {
	handler := &Handler{
		service:    service,
		jwtManager: jwtManager,
	}
	router.HandleFunc("POST /register", handler.Register)
	router.HandleFunc("POST /login", handler.Login)
	router.Handle("DELETE /user", auth.Middleware(handler.jwtManager, http.HandlerFunc(handler.Delete)))
}

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var payload RegisterRequest

	err := decodeJSON[RegisterRequest](r, &payload)
	if err != nil {
		writeJSON(w, err, http.StatusBadRequest)
		return
	}

	err = validateStruct[RegisterRequest](payload)
	if err != nil {
		writeJSON(w, err, http.StatusBadRequest)
		return
	}

	u, err := handler.service.Register(r.Context(), payload.Email, payload.Password)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := RegisterResponse{
		ID:    u.ID,
		Email: u.Email,
	}

	writeJSON(w, data, http.StatusCreated)
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var payload LoginRequest

	err := decodeJSON[LoginRequest](r, &payload)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateStruct[LoginRequest](payload)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := handler.service.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := handler.jwtManager.Create(auth.JWTData{UserID: user.ID})
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
	}

	data := LoginResponse{
		Token: token,
	}

	writeJSON(w, data, http.StatusOK)

}

func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, "", http.StatusUnauthorized)
		return
	}

	err := handler.service.Delete(r.Context(), userID)
	if err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, "user deleted", http.StatusOK)
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
