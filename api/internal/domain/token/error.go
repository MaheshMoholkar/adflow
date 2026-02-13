package token

import "errors"

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token has expired")
	ErrTokenRevoked  = errors.New("token has been revoked")
)
