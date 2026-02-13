package auth

import "context"

// Service defines the interface for authentication business logic
type Service interface {
	// Register creates a new user with phone and password
	Register(ctx context.Context, req RegisterRequest) (*TokenResponse, error)

	// Login authenticates a user with phone and password
	Login(ctx context.Context, req LoginRequest) (*TokenResponse, error)

	// VerifyToken verifies a JWT token and returns the claims
	VerifyToken(ctx context.Context, tokenString string) (*AuthClaims, error)
}
