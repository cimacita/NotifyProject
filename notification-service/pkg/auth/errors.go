package auth

import "errors"

var (
	ErrSignMethod        = errors.New("unexpected signing method")
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
	ErrMissingUserID     = errors.New("missing user_id in token")
)
