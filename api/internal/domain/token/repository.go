package token

import (
	"context"
	"time"
)

// Repository defines the interface for token data access
type Repository interface {
	CreateRefreshToken(ctx context.Context, data TokenCreate) (*Token, error)
	GetRefreshTokenByToken(ctx context.Context, token string) (*Token, error)
	RevokeToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID int64) error
	UpdateTokenLastUsed(ctx context.Context, id int64) error
	DeleteExpiredTokens(ctx context.Context, before time.Time) error
	StoreTokenID(ctx context.Context, tokenID string, userID int64, expiresAt time.Time, tokenType string) error
	IsTokenRevoked(ctx context.Context, tokenID string) (bool, error)
	RevokeTokenByID(ctx context.Context, tokenID string) error
	RevokeAllTokensByType(ctx context.Context, userID int64, tokenType string) error
}
