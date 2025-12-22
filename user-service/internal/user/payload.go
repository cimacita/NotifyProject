package user

import "github.com/google/uuid"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
