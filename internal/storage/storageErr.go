package storage

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExist   = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCodeBlocked        = errors.New("too many attempts, try again later")
	ErrCodeInvalid        = errors.New("invalid code")
	ErrCodeNotFound       = errors.New("code not found or expired")
)
