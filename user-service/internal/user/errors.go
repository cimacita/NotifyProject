package user

import "errors"

var (
	ErrUserExists       = errors.New("user already exists")
	ErrWrongCredentials = errors.New("wrong credentials")
	ErrUserNotFound     = errors.New("user not found")
)
