package user

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrPhoneTaken   = errors.New("phone number already registered")
)
