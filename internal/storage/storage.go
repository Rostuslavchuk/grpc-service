package storage

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user no found")
	ErrPasswordIncorect   = errors.New("incorect password")
	ErrAppNotFound        = errors.New("app not found")
)
