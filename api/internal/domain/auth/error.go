package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid phone or password")
	ErrPhoneTaken         = errors.New("phone number already registered")
	ErrExpiredToken       = errors.New("token has expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrUnauthorized       = errors.New("unauthorized")
)
